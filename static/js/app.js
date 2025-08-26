(async function () {
  let i18n = {};
  let lang = 'en';
  let subformsCache = null;
  let categoriesCache = [];

  const $ = (id) => document.getElementById(id);

  const t = (path) => {
    const parts = path.split('.');
    let cur = i18n;
    for (const p of parts) {
      if (cur && typeof cur === 'object' && p in cur) cur = cur[p];
      else return path;
    }
    return typeof cur === 'string' ? cur : path;
  };

  async function loadLang(l) {
    const res = await fetch(`/api/translations/${l}`);
    if (!res.ok) throw new Error('i18n load failed');
    i18n = await res.json();
    $('app-title').textContent = i18n.app_title || 'Insurance Quote Dashboard';
    $('welcome').textContent = i18n.welcome || '';
    $('formTitle').textContent = 'Form';
  }

  async function loadCategories() {
    const res = await fetch('/api/category/list');
    categoriesCache = await res.json();
  }

  // Build tab bar (one tab per category, ordered)
  function buildTabs(activeId) {
    const wrap = $('catTabs');
    wrap.innerHTML = '';
    categoriesCache.forEach((c, idx) => {
      const btn = document.createElement('button');
      btn.className = 'tab-btn';
      btn.setAttribute('role', 'tab');
      btn.setAttribute('aria-selected', activeId ? String(c.id === activeId) : String(idx === 0));
      btn.setAttribute('aria-controls', 'formWrap');
      btn.dataset.cat = c.id;
      btn.title = t(c.description) || '';
      btn.innerHTML = `<span class="tab-ico">${c.icon || ''}</span><span class="tab-txt">${t(c.title) || c.title}</span>`;
      btn.onclick = () => {
        setActiveTab(c.id);
        openCategory(c.id);
      };
      wrap.appendChild(btn);
    });
    // visual active state
    setActiveTab(activeId || (categoriesCache[0] && categoriesCache[0].id));
  }

  function setActiveTab(catId) {
    const buttons = document.querySelectorAll('#catTabs .tab-btn');
    buttons.forEach(b => {
      const isActive = b.dataset.cat === catId;
      b.classList.toggle('active', isActive);
      b.setAttribute('aria-selected', String(isActive));
      b.tabIndex = isActive ? 0 : -1;
    });
  }

  // Fetch category schema + render form
  async function openCategory(cat) {
    const res = await fetch(`/api/category/${cat}`);
    const data = await res.json();
    subformsCache = data.subforms || {};
    renderForm(cat, data.fields);
    $('formTitle').textContent = t(`categories.${cat}.title`) || cat;
  }

  // ===== Form rendering =====
  function renderForm(cat, fields) {
    const wrap = $('formWrap');
    wrap.innerHTML = '';
    const form = document.createElement('form');
    form.id = 'catForm';

    fields.forEach(f => {
      form.appendChild(buildField(f, cat));
      if (cat === 'vehicle' && f.property === 'registrationNumber') {
        form.appendChild(dvlaControls());
      }
      if (cat === 'vehicle' && f.property === 'hasModifications') {
        const area = document.createElement('div');
        area.id = 'modsArea';
        form.appendChild(area);
      }
    });

    if (cat === 'payment') addPaymentExtras(form);

    const saveBtn = document.createElement('button');
    saveBtn.type = 'button';
    saveBtn.textContent = i18n.save_continue || 'Save & Continue';
    saveBtn.onclick = () => saveForm(cat, form);
    saveBtn.className = 'btn';
    form.appendChild(document.createElement('hr'));
    form.appendChild(saveBtn);

    wrap.appendChild(form);
  }

  // Helpers to build inputs
  function labelEl(text, forId) {
    const l = document.createElement('label');
    l.textContent = text;
    if (forId) l.setAttribute('for', forId);
    return l;
  }

  function buildField(f, category) {
    const row = document.createElement('div');
    row.className = 'row';

    const id = `f_${f.property}`;
    const lbl = labelEl(t(f.label) || f.label, id);
    if (f.required) {
      const star = document.createElement('span'); star.textContent = ' *'; lbl.appendChild(star);
    }
    row.appendChild(lbl);

    let input;
    if (f.type === 'select') {
      input = document.createElement('select');
      input.id = id;
      input.dataset.prop = f.property;
      
      // Special handling for relationship field in drivers category
      if (f.property === 'relationshipToMainDriver' && category === 'drivers') {
        // Get current driver index from the form context
        const driverIndex = getCurrentDriverIndex();
        const isMainDriver = driverIndex === 0; // First driver is always main driver
        
        (f.options || []).forEach(o => {
          // Only show SELF option for main driver
          if (o.value === 'SELF' && !isMainDriver) {
            return; // Skip SELF option for non-main drivers
          }
          const opt = document.createElement('option');
          opt.value = o.value;
          opt.textContent = t(o.label) || o.label;
          input.appendChild(opt);
        });
      } else {
        (f.options || []).forEach(o => {
          const opt = document.createElement('option');
          opt.value = o.value;
          opt.textContent = t(o.label) || o.label;
          input.appendChild(opt);
        });
      }
    } else if (f.type === 'radio') {
      input = document.createElement('div');
      input.id = id;
      input.dataset.prop = f.property;
      (f.options || []).forEach(o => {
        const rId = `${id}_${o.value}`;
        const r = document.createElement('input');
        r.type = 'radio';
        r.name = id;
        r.value = o.value;
        r.id = rId;
        const rl = labelEl(t(o.label) || o.label, rId);
        input.appendChild(r);
        input.appendChild(rl);
      });
      if (f.property === 'hasModifications') {
        input.addEventListener('change', () => {
          const val = selectedRadioValue(id);
          const area = $('modsArea');
          if (area) {
            area.innerHTML = '';
            if (val === 'YES') renderModifications(area);
          }
        });
      }
    } else {
      input = document.createElement('input');
      input.id = id;
      input.dataset.prop = f.property;
      input.type = f.type === 'date' ? 'date'
               : f.type === 'number' ? 'number'
               : f.type === 'email' ? 'email'
               : f.type === 'tel' ? 'tel' : 'text';

      if (f.property === 'phone') input.addEventListener('input', maskUKMobile);
      if (f.property === 'postcode') input.addEventListener('input', upcaseNoSpacesEndTrim);
      if (f.property === 'licenceNumber') input.addEventListener('input', upcaseNoSpaces);
      if (f.property === 'registrationNumber') input.addEventListener('input', upcaseTrim);
    }

    if (f.helpText) {
      const help = document.createElement('div');
      help.className = 'help';
      help.textContent = t(f.helpText) || '';
      row.appendChild(help);
    }

    row.appendChild(input);
    return row;
  }

  function selectedRadioValue(id) {
    const radios = document.querySelectorAll(`#${id} input[type="radio"]`);
    for (const r of radios) if (r.checked) return r.value;
    return '';
  }

  function getCurrentDriverIndex() {
    // This function should return the current driver index being edited
    // For now, we'll assume it's the first driver (index 0) for the main driver
    // In a more complex implementation, this would be determined by the current form context
    return 0;
  }

  async function validateDriverRelationship(driverIndex, fields) {
    try {
      const response = await fetch('/api/drivers/validate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Session-ID': getSessionID()
        },
        body: JSON.stringify({
          driverIndex: driverIndex,
          fields: fields
        })
      });
      
      if (response.ok) {
        const result = await response.json();
        return result;
      }
    } catch (error) {
      console.error('Driver validation error:', error);
    }
    return { valid: true, errors: [] };
  }

  function getSessionID() {
    // Get session ID from cookie or other storage
    // This is a simplified implementation
    return document.cookie.split('; ').find(row => row.startsWith('quote-session='))?.split('=')[1] || '';
  }

  // ===== Modifications subforms =====
  function renderModifications(area) {
    const mods = subformsCache?.modifications;
    if (!mods) return;

    const selRow = document.createElement('div');
    selRow.className = 'row';
    const lbl = labelEl(t('fields.modificationType') || 'Modification Category');
    selRow.appendChild(lbl);

    const select = document.createElement('select');
    select.multiple = true;
    select.size = 6;
    select.id = 'modificationType';
    (mods.fields?.[0]?.options || []).forEach(o => {
      const opt = document.createElement('option');
      opt.value = o.value;
      opt.textContent = o.label;
      select.appendChild(opt);
    });
    selRow.appendChild(select);
    area.appendChild(selRow);

    const subwrap = document.createElement('div');
    subwrap.id = 'modsSubforms';
    area.appendChild(subwrap);

    select.addEventListener('change', () => {
      subwrap.innerHTML = '';
      const chosen = Array.from(select.selectedOptions).map(o => o.value);
      for (const key in mods.subforms) {
        const sub = mods.subforms[key];
        if (chosen.includes(sub.triggerValue)) {
          subwrap.appendChild(renderSubform(sub));
        }
      }
    });
  }

  function renderSubform(sub) {
    const box = document.createElement('div');
    box.className = 'card';
    const ttl = document.createElement('div');
    ttl.className = 'title';
    ttl.textContent = sub.title;
    box.appendChild(ttl);

    (sub.fields || []).forEach(f => {
      const fieldDef = { ...f, label: f.label || f.property, type: f.type || 'text' };
      box.appendChild(buildField(fieldDef));
      if (f.type === 'number') {
        const el = box.querySelector(`#f_${f.property}`);
        if (el) {
          if (typeof f.min === 'number') el.min = f.min;
          if (typeof f.max === 'number') el.max = f.max;
        }
      }
    });
    return box;
  }

  // ===== DVLA controls =====
  function dvlaControls() {
    const row = document.createElement('div');
    row.className = 'row-inline';
    const liveId = 'dvla_live';
    row.innerHTML = `
      <label><input id="${liveId}" type="checkbox"/> Live DVLA (Node-RED)</label>
      <button type="button" id="dvlaBtn">Lookup</button>
      <span id="dvlaStatus" class="muted"></span>
    `;
    setTimeout(() => {
      $('dvlaBtn').onclick = async () => {
        const reg = $('f_registrationNumber').value || '';
        const live = $(liveId).checked ? 'true' : 'false';
        if (!reg) return;
        $('dvlaStatus').textContent = '…';
        const q = new URLSearchParams({ reg, live }).toString();
        const res = await fetch(`/api/dvla/lookup?${q}`);
        if (res.ok) {
          const j = await res.json();
          applyDVLAResult(j);
          $('dvlaStatus').textContent = '✓';
        } else {
          $('dvlaStatus').textContent = 'lookup failed';
        }
      };
    });
    return row;
  }

  function applyDVLAResult(j) {
    const set = (prop, val) => {
      const el = $(`f_${prop}`);
      if (el && !el.value) el.value = val;
    };
    set('vehicleMake', j.make);
    set('vehicleModel', j.model);
    set('vehicleYear', j.year);
  }

  // ===== Payment extras (card/direct debit) =====
  function addPaymentExtras(form) {
    const holder = document.createElement('div');
    holder.id = 'paymentExtras';
    form.appendChild(holder);

    const methodSel = form.querySelector('#f_paymentMethod');
    if (methodSel) methodSel.addEventListener('change', () => renderPaymentExtras(holder, methodSel.value));
    if (methodSel) renderPaymentExtras(holder, methodSel.value || '');
  }

  function renderPaymentExtras(holder, method) {
    holder.innerHTML = '';
    if (method === 'CREDIT_CARD' || method === 'DEBIT_CARD') {
      holder.appendChild(rowInput('Card Number', 'cardNumber', { oninput: maskCardNumber }));
      holder.appendChild(rowInput('Expiry (MM/YY)', 'cardExpiry', { oninput: maskExpiry }));
      holder.appendChild(rowInput('CVC', 'cardCVC', { type: 'tel', maxlength: 4 }));
    } else if (method === 'DIRECT_DEBIT') {
      holder.appendChild(rowInput('Account Name', 'accountName'));
      holder.appendChild(rowInput('Sort Code (##-##-##)', 'sortCode', { oninput: maskSortCode, maxlength: 8 }));
      holder.appendChild(rowInput('Account Number (8 digits)', 'accountNumber', { oninput: digitsOnly, maxlength: 8 }));
    }
  }

  function rowInput(label, prop, opts = {}) {
    const row = document.createElement('div');
    row.className = 'row';
    row.appendChild(labelEl(label, `f_${prop}`));
    const input = document.createElement('input');
    input.id = `f_${prop}`;
    input.dataset.prop = prop;
    input.type = opts.type || 'text';
    if (opts.maxlength) input.maxLength = opts.maxlength;
    if (opts.oninput) input.addEventListener('input', opts.oninput);
    row.appendChild(input);
    return row;
  }

  // ===== Save form =====
  async function saveForm(cat, form) {
    const fields = {};
    form.querySelectorAll('[data-prop]').forEach(el => {
      if (el.tagName === 'DIV' && el.querySelector('input[type="radio"]')) {
        const radios = el.querySelectorAll('input[type="radio"]');
        let v = '';
        radios.forEach(r => { if (r.checked) v = r.value; });
        fields[el.dataset.prop] = v;
      } else if (el.tagName === 'SELECT' && el.multiple) {
        fields[el.dataset.prop] = Array.from(el.selectedOptions).map(o => o.value);
      } else {
        fields[el.dataset.prop] = el.value;
      }
    });

    // Special validation for drivers category
    if (cat === 'drivers') {
      const driverIndex = getCurrentDriverIndex();
      const driverValidation = await validateDriverRelationship(driverIndex, fields);
      if (!driverValidation.valid) {
        setStatus(`${i18n.validation_error || 'Please correct errors'}: ${driverValidation.errors.map(e=>e.field).join(', ')}`);
        return;
      }
    }

    if (cat === 'payment') {
      const method = fields.paymentMethod;
      if (method === 'CREDIT_CARD' || method === 'DEBIT_CARD') {
        const num = ($('f_cardNumber')?.value || '').replace(/\s+/g, '');
        if (!luhn(num)) return setStatus('Invalid card number');
        const exp = $('f_cardExpiry')?.value || '';
        if (!/^(0[1-9]|1[0-2])\/\d{2}$/.test(exp)) return setStatus('Invalid expiry (MM/YY)');
        const cvc = $('f_cardCVC')?.value || '';
        if (!/^\d{3,4}$/.test(cvc)) return setStatus('Invalid CVC');
        fields.cardNumber = num;
        fields.cardExpiry = exp;
        fields.cardCVC = cvc;
      } else if (method === 'DIRECT_DEBIT') {
        const sc = (($('f_sortCode')?.value || '').match(/\d/g) || []).join('');
        const an = $('f_accountNumber')?.value || '';
        if (sc.length !== 6) return setStatus('Invalid sort code');
        if (!/^\d{8}$/.test(an)) return setStatus('Invalid account number');
        fields.sortCode = sc;
        fields.accountNumber = an;
        fields.accountName = $('f_accountName')?.value || '';
      }
    }

    let res = await fetch('/api/validate', {
      method: 'POST', headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ category: cat, fields })
    });
    const vr = await res.json();
    if (!vr.valid) {
      setStatus(`${i18n.validation_error || 'Please correct errors'}: ${vr.errors.map(e=>e.field).join(', ')}`);
      return;
    }

    res = await fetch('/api/save', {
      method: 'POST', headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ category: cat, fields })
    });
    if (res.ok) setStatus('Saved.');
  }

  function setStatus(msg) { $('status').textContent = msg; }

  // ===== masks & helpers =====
  function maskUKMobile(e) {
    e.target.value = e.target.value.replace(/[^\d]/g, '').slice(0, 11);
  }
  function upcaseNoSpaces(e){ e.target.value = e.target.value.toUpperCase().replace(/\s+/g,''); }
  function upcaseNoSpacesEndTrim(e){ e.target.value = e.target.value.toUpperCase().replace(/[^A-Z0-9 ]/g,'').replace(/\s+/g,' ').trim(); }
  function upcaseTrim(e){ e.target.value = e.target.value.toUpperCase().replace(/\s+/g,'').trim(); }
  function maskCardNumber(e){
    const d = e.target.value.replace(/\D/g,'').slice(0,19);
    e.target.value = d.replace(/(.{4})/g,'$1 ').trim();
  }
  function maskExpiry(e){
    let v = e.target.value.replace(/\D/g,'').slice(0,4);
    if (v.length >= 3) v = v.slice(0,2) + '/' + v.slice(2);
    e.target.value = v;
  }
  function maskSortCode(e){
    const d = e.target.value.replace(/\D/g,'').slice(0,6);
    e.target.value = d.replace(/(\d{2})(?=\d)/g,'$1-');
  }
  function digitsOnly(e){ e.target.value = e.target.value.replace(/\D/g,''); }
  function luhn(num){
    if (!/^\d{12,19}$/.test(num)) return false;
    let sum=0, dbl=false;
    for (let i=num.length-1;i>=0;i--){
      let n = +num[i];
      if (dbl){ n*=2; if(n>9)n-=9; }
      sum+=n; dbl=!dbl;
    }
    return sum%10===0;
  }

  // Language + session controls
  $('langSel').addEventListener('change', async (e) => {
    lang = e.target.value;
    await loadLang(lang);
    // rebuild tabs with translated titles
    buildTabs(document.querySelector('#catTabs .tab-btn.active')?.dataset.cat);
  });

  $('newSession').addEventListener('click', async () => {
    await fetch('/api/session/new', { method:'POST' });
    setStatus('Started a new session.');
  });

  $('exportSession').addEventListener('click', async () => {
    const res = await fetch('/api/session/export');
    if (!res.ok) return;
    const blob = await res.blob();
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'quote-session.json';
    a.click();
  });

  $('importFile').addEventListener('change', async (e) => {
    const file = e.target.files[0];
    if (!file) return;
    const text = await file.text();
    await fetch('/api/session/import', { method:'POST', headers:{'Content-Type':'application/json'}, body:text });
    setStatus('Imported session.');
  });

  // Boot
  await loadLang(lang);
  await loadCategories();
  buildTabs();                     // build the tab bar
  if (categoriesCache[0]) {
    await openCategory(categoriesCache[0].id);  // open first tab by default
  }
})();

