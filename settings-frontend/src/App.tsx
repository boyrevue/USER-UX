import React, { useState } from 'react';
import { 
  Building2, 
  CreditCard as CreditCardIcon, 
  MessageSquare, 
  Shield, 
  Plus,
  X,
  Camera,
  Upload,
  HelpCircle,
  Smartphone,
  Mic,
  Lock,
  Mail,
  Settings,
  Eye,
  EyeOff,
  Download,
  Trash2
} from 'lucide-react';
import { 
  Card, 
  Button,
  Badge,
  Modal,
  TextInput,
  Label,
  Select,
  Textarea,
  Tabs,
  TabItem,
  ToggleSwitch
} from 'flowbite-react';

interface BankAccount {
  id: number;
  bankName: string;
  accountNumber: string;
  sortCode: string;
  accountHolder: string;
  openBankingEnabled: boolean;
  truelayerClientId: string;
  truelayerClientSecret: string;
  truelayerRedirectUri: string;
  truelayerEnvironment: string;
  truelayerScopes: string[];
}

interface CreditCard {
  id: number;
  cardNumber: string;
  expiryDate: string;
  cvv: string;
  cardholderName: string;
  billingAddress: string;
  cardType: string;
  provider: string;
  openBankingEnabled: boolean;
  truelayerClientId: string;
  truelayerClientSecret: string;
  truelayerRedirectUri: string;
  truelayerEnvironment: string;
  truelayerScopes: string[];
}

interface EmailConfig {
  provider: string;
  emailAddress: string;
  oauth2Enabled: boolean;
  oauth2ClientId: string;
  oauth2ClientSecret: string;
  oauth2RefreshToken: string;
  oauth2AccessToken: string;
  password: string;
  autoSyncEnabled: boolean;
  syncFrequency: string;
  imapServer: string;
  imapPort: string;
  smtpServer: string;
  smtpPort: string;
  useSSL: boolean;
  useTLS: boolean;
}

interface CommunicationChannel {
  id: number;
  type: 'email' | 'sms' | 'voice' | 'secure';
  enabled: boolean;
  config: any;
}

function App() {
  const [activeTab, setActiveTab] = useState(0);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([
    {
      id: 1,
      bankName: '',
      accountNumber: '',
      sortCode: '',
      accountHolder: '',
      openBankingEnabled: false,
      truelayerClientId: '',
      truelayerClientSecret: '',
      truelayerRedirectUri: '',
      truelayerEnvironment: 'sandbox',
      truelayerScopes: []
    }
  ]);
  const [creditCards, setCreditCards] = useState<CreditCard[]>([
    {
      id: 1,
      cardNumber: '',
      expiryDate: '',
      cvv: '',
      cardholderName: '',
      billingAddress: '',
      cardType: 'visa',
      provider: 'visa',
      openBankingEnabled: false,
      truelayerClientId: '',
      truelayerClientSecret: '',
      truelayerRedirectUri: '',
      truelayerEnvironment: 'sandbox',
      truelayerScopes: []
    }
  ]);
  const [emailConfig, setEmailConfig] = useState<EmailConfig>({
    provider: '',
    emailAddress: '',
    oauth2Enabled: false,
    oauth2ClientId: '',
    oauth2ClientSecret: '',
    oauth2RefreshToken: '',
    oauth2AccessToken: '',
    password: '',
    autoSyncEnabled: false,
    syncFrequency: '5min',
    imapServer: '',
    imapPort: '993',
    smtpServer: '',
    smtpPort: '587',
    useSSL: true,
    useTLS: true
  });
  const [communicationChannels] = useState<CommunicationChannel[]>([
    { id: 1, type: 'email', enabled: true, config: {} },
    { id: 2, type: 'sms', enabled: false, config: {} },
    { id: 3, type: 'voice', enabled: false, config: {} },
    { id: 4, type: 'secure', enabled: false, config: {} }
  ]);

  const [isUploadOpen, setIsUploadOpen] = useState(false);
  const [isHelpOpen, setIsHelpOpen] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showOAuthSecret, setShowOAuthSecret] = useState(false);

  const addBankAccount = () => {
    const newAccount: BankAccount = {
      id: bankAccounts.length + 1,
      bankName: '',
      accountNumber: '',
      sortCode: '',
      accountHolder: '',
      openBankingEnabled: false,
      truelayerClientId: '',
      truelayerClientSecret: '',
      truelayerRedirectUri: '',
      truelayerEnvironment: 'sandbox',
      truelayerScopes: []
    };
    setBankAccounts([...bankAccounts, newAccount]);
  };

  const addCreditCard = () => {
    const newCard: CreditCard = {
      id: creditCards.length + 1,
      cardNumber: '',
      expiryDate: '',
      cvv: '',
      cardholderName: '',
      billingAddress: '',
      cardType: 'visa',
      provider: 'visa',
      openBankingEnabled: false,
      truelayerClientId: '',
      truelayerClientSecret: '',
      truelayerRedirectUri: '',
      truelayerEnvironment: 'sandbox',
      truelayerScopes: []
    };
    setCreditCards([...creditCards, newCard]);
  };

  const removeBankAccount = (id: number) => {
    setBankAccounts(bankAccounts.filter(account => account.id !== id));
  };

  const removeCreditCard = (id: number) => {
    setCreditCards(creditCards.filter(card => card.id !== id));
  };

  const updateBankAccount = (id: number, field: keyof BankAccount, value: any) => {
    setBankAccounts(bankAccounts.map(account => 
      account.id === id ? { ...account, [field]: value } : account
    ));
  };

  const updateCreditCard = (id: number, field: keyof CreditCard, value: any) => {
    setCreditCards(creditCards.map(card => 
      card.id === id ? { ...card, [field]: value } : card
    ));
  };

  const updateEmailConfig = (field: keyof EmailConfig, value: any) => {
    setEmailConfig({ ...emailConfig, [field]: value });
  };

  const getStatusBadge = (isRequired: boolean, hasValue: boolean) => {
    if (isRequired && !hasValue) {
      return <Badge color="failure" className="status-badge status-missing">Missing</Badge>;
    } else if (isRequired) {
      return <Badge color="success" className="status-badge status-required">REQ</Badge>;
    } else {
      return <Badge color="warning" className="status-badge status-optional">Optional</Badge>;
    }
  };

  const formatCardNumber = (value: string) => {
    const v = value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
    const matches = v.match(/\d{4,16}/g);
    const match = matches && matches[0] || '';
    const parts = [];
    for (let i = 0, len = match.length; i < len; i += 4) {
      parts.push(match.substring(i, i + 4));
    }
    if (parts.length) {
      return parts.join(' ');
    } else {
      return v;
    }
  };

  const getCardType = (cardNumber: string) => {
    const cleanNumber = cardNumber.replace(/\s/g, '');
    if (/^4/.test(cleanNumber)) return 'visa';
    if (/^5[1-5]/.test(cleanNumber)) return 'mastercard';
    if (/^3[47]/.test(cleanNumber)) return 'amex';
    if (/^6/.test(cleanNumber)) return 'discover';
    return 'unknown';
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* Header */}
      <div className="border-b border-gray-700 bg-gray-800">
        <div className="flex items-center justify-between p-4">
          <div>
            <h1 className="text-xl font-semibold">Settings Configuration</h1>
            <p className="text-gray-400 text-sm">Manage your application settings</p>
          </div>
          <div className="flex gap-2">
            <Button color="gray" onClick={() => setIsHelpOpen(true)}>
              <HelpCircle className="w-4 h-4 mr-2" />
              Get Help
            </Button>
            <Button color="gray" onClick={() => setIsUploadOpen(true)}>
              <Upload className="w-4 h-4 mr-2" />
              Upload Documents
            </Button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex h-[calc(100vh-80px)]">
        {/* Sidebar */}
        <div className="w-80 bg-gray-800 border-r border-gray-700">
          <Tabs aria-label="Settings tabs" onActiveTabChange={(tab) => setActiveTab(tab)}>
            <TabItem active title={
              <div className="flex items-center gap-2">
                <Building2 className="w-4 h-4" />
                Bank Accounts
                <Badge color="gray" className="ml-auto">{bankAccounts.length}</Badge>
              </div>
            }>
              <div className="p-6 space-y-4">
                <div className="flex items-center justify-between">
                  <h2 className="text-lg font-semibold">Bank Accounts</h2>
                  <Button size="sm" onClick={addBankAccount}>
                    <Plus className="w-4 h-4 mr-2" />
                    Add Bank
                  </Button>
                </div>
                
                {bankAccounts.map((account, index) => (
                  <Card key={account.id} className="bg-gray-700 border-gray-600">
                    <div className="flex items-center justify-between pb-3">
                      <h3 className="text-sm font-semibold">Bank Account {index + 1}</h3>
                      {bankAccounts.length > 1 && (
                        <Button 
                          size="sm" 
                          color="failure" 
                          onClick={() => removeBankAccount(account.id)}
                        >
                          <X className="w-4 h-4" />
                        </Button>
                      )}
                    </div>
                    <div className="space-y-4">
                      <div>
                        <Label className="flex items-center justify-between">
                          Bank Name
                          {getStatusBadge(true, !!account.bankName)}
                        </Label>
                        <Select 
                          value={account.bankName} 
                          onChange={(e) => updateBankAccount(account.id, 'bankName', e.target.value)}
                        >
                          <option value="">Select Bank</option>
                          <option value="barclays">Barclays</option>
                          <option value="hsbc">HSBC</option>
                          <option value="lloyds">Lloyds Bank</option>
                          <option value="natwest">NatWest</option>
                          <option value="santander">Santander</option>
                          <option value="rbs">Royal Bank of Scotland</option>
                          <option value="nationwide">Nationwide</option>
                          <option value="tsb">TSB</option>
                          <option value="coop">Co-operative Bank</option>
                          <option value="metrobank">Metro Bank</option>
                        </Select>
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Account Number
                          {getStatusBadge(true, !!account.accountNumber)}
                        </Label>
                        <TextInput
                          value={account.accountNumber}
                          onChange={(e) => updateBankAccount(account.id, 'accountNumber', e.target.value)}
                          placeholder="12345678"
                          maxLength={8}
                          className="bg-gray-600 border-gray-500 text-white"
                        />
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Sort Code
                          {getStatusBadge(true, !!account.sortCode)}
                        </Label>
                        <TextInput
                          value={account.sortCode}
                          onChange={(e) => updateBankAccount(account.id, 'sortCode', e.target.value)}
                          placeholder="12-34-56"
                          maxLength={8}
                          className="bg-gray-600 border-gray-500 text-white"
                        />
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Account Holder
                          {getStatusBadge(true, !!account.accountHolder)}
                        </Label>
                        <TextInput
                          value={account.accountHolder}
                          onChange={(e) => updateBankAccount(account.id, 'accountHolder', e.target.value)}
                          placeholder="John Doe"
                          className="bg-gray-600 border-gray-500 text-white"
                        />
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Open Banking
                          {getStatusBadge(false, true)}
                        </Label>
                        <ToggleSwitch
                          checked={account.openBankingEnabled}
                          onChange={(checked) => updateBankAccount(account.id, 'openBankingEnabled', checked)}
                        />
                      </div>

                      {account.openBankingEnabled && (
                        <div className="space-y-4 p-4 bg-gray-600 rounded-lg">
                          <h4 className="font-medium text-sm">TrueLayer Open Banking Configuration</h4>
                          
                          <div>
                            <Label className="flex items-center justify-between">
                              Client ID
                              {getStatusBadge(true, !!account.truelayerClientId)}
                            </Label>
                            <TextInput
                              value={account.truelayerClientId}
                              onChange={(e) => updateBankAccount(account.id, 'truelayerClientId', e.target.value)}
                              placeholder="your-truelayer-client-id"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>

                          <div>
                            <Label className="flex items-center justify-between">
                              Client Secret
                              {getStatusBadge(true, !!account.truelayerClientSecret)}
                            </Label>
                            <div className="relative">
                              <TextInput
                                type={showOAuthSecret ? "text" : "password"}
                                value={account.truelayerClientSecret}
                                onChange={(e) => updateBankAccount(account.id, 'truelayerClientSecret', e.target.value)}
                                placeholder="your-truelayer-client-secret"
                                className="bg-gray-500 border-gray-400 text-white pr-10"
                              />
                              <button
                                type="button"
                                className="absolute right-2 top-1/2 transform -translate-y-1/2"
                                onClick={() => setShowOAuthSecret(!showOAuthSecret)}
                              >
                                {showOAuthSecret ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                              </button>
                            </div>
                          </div>

                          <div>
                            <Label className="flex items-center justify-between">
                              Redirect URI
                              {getStatusBadge(true, !!account.truelayerRedirectUri)}
                            </Label>
                            <TextInput
                              value={account.truelayerRedirectUri}
                              onChange={(e) => updateBankAccount(account.id, 'truelayerRedirectUri', e.target.value)}
                              placeholder="https://yourapp.com/oauth/callback"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>

                          <div>
                            <Label className="flex items-center justify-between">
                              Environment
                              {getStatusBadge(false, true)}
                            </Label>
                            <Select 
                              value={account.truelayerEnvironment} 
                              onChange={(e) => updateBankAccount(account.id, 'truelayerEnvironment', e.target.value)}
                            >
                              <option value="sandbox">Sandbox (Testing)</option>
                              <option value="live">Live (Production)</option>
                            </Select>
                          </div>
                        </div>
                      )}
                    </div>
                  </Card>
                ))}
              </div>
            </TabItem>

            <TabItem title={
              <div className="flex items-center gap-2">
                <CreditCardIcon className="w-4 h-4" />
                Credit Cards
                <Badge color="gray" className="ml-auto">{creditCards.length}</Badge>
              </div>
            }>
              <div className="p-6 space-y-4">
                <div className="flex items-center justify-between">
                  <h2 className="text-lg font-semibold">Credit Cards</h2>
                  <Button size="sm" onClick={addCreditCard}>
                    <Plus className="w-4 h-4 mr-2" />
                    Add Card
                  </Button>
                </div>
                
                {creditCards.map((card, index) => (
                  <Card key={card.id} className="bg-gray-700 border-gray-600">
                    <div className="flex items-center justify-between pb-3">
                      <h3 className="text-sm font-semibold">Credit Card {index + 1}</h3>
                      <div className="flex gap-2">
                        <Button size="sm" color="gray">
                          <Camera className="w-4 h-4" />
                        </Button>
                        {creditCards.length > 1 && (
                          <Button 
                            size="sm" 
                            color="failure" 
                            onClick={() => removeCreditCard(card.id)}
                          >
                            <X className="w-4 h-4" />
                          </Button>
                        )}
                      </div>
                    </div>
                    
                    {/* Credit Card Visual Display */}
                    <div className="mb-4 p-4 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg text-white">
                      <div className="flex justify-between items-start mb-4">
                        <div className="text-sm opacity-80">Credit Card</div>
                        <div className="text-sm opacity-80">{card.cardType.toUpperCase()}</div>
                      </div>
                      <div className="text-lg font-mono mb-4">
                        {card.cardNumber ? formatCardNumber(card.cardNumber) : '•••• •••• •••• ••••'}
                      </div>
                      <div className="flex justify-between items-center">
                        <div>
                          <div className="text-xs opacity-80">Cardholder</div>
                          <div className="text-sm">{card.cardholderName || 'YOUR NAME'}</div>
                        </div>
                        <div>
                          <div className="text-xs opacity-80">Expires</div>
                          <div className="text-sm">{card.expiryDate || 'MM/YY'}</div>
                        </div>
                      </div>
                    </div>

                    <div className="space-y-4">
                      <div>
                        <Label className="flex items-center justify-between">
                          Card Number
                          {getStatusBadge(true, !!card.cardNumber)}
                        </Label>
                        <TextInput
                          value={card.cardNumber}
                          onChange={(e) => updateCreditCard(card.id, 'cardNumber', formatCardNumber(e.target.value))}
                          placeholder="1234 5678 9012 3456"
                          maxLength={19}
                          className="bg-gray-600 border-gray-500 text-white"
                        />
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label className="flex items-center justify-between">
                            Expiry Date
                            {getStatusBadge(true, !!card.expiryDate)}
                          </Label>
                          <TextInput
                            value={card.expiryDate}
                            onChange={(e) => updateCreditCard(card.id, 'expiryDate', e.target.value)}
                            placeholder="MM/YY"
                            maxLength={5}
                            className="bg-gray-600 border-gray-500 text-white"
                          />
                        </div>
                        <div>
                          <Label className="flex items-center justify-between">
                            CVV
                            {getStatusBadge(true, !!card.cvv)}
                          </Label>
                          <TextInput
                            value={card.cvv}
                            onChange={(e) => updateCreditCard(card.id, 'cvv', e.target.value)}
                            placeholder="123"
                            maxLength={4}
                            className="bg-gray-600 border-gray-500 text-white"
                          />
                        </div>
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Cardholder Name
                          {getStatusBadge(true, !!card.cardholderName)}
                        </Label>
                        <TextInput
                          value={card.cardholderName}
                          onChange={(e) => updateCreditCard(card.id, 'cardholderName', e.target.value)}
                          placeholder="John Doe"
                          className="bg-gray-600 border-gray-500 text-white"
                        />
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Billing Address
                          {getStatusBadge(false, true)}
                        </Label>
                        <Textarea
                          value={card.billingAddress}
                          onChange={(e) => updateCreditCard(card.id, 'billingAddress', e.target.value)}
                          placeholder="Enter billing address"
                          className="bg-gray-600 border-gray-500 text-white"
                        />
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Card Type
                          {getStatusBadge(false, true)}
                        </Label>
                        <Select 
                          value={card.cardType} 
                          onChange={(e) => updateCreditCard(card.id, 'cardType', e.target.value)}
                        >
                          <option value="visa">Visa</option>
                          <option value="mastercard">Mastercard</option>
                          <option value="amex">American Express</option>
                          <option value="discover">Discover</option>
                        </Select>
                      </div>

                      <div>
                        <Label className="flex items-center justify-between">
                          Open Banking
                          {getStatusBadge(false, true)}
                        </Label>
                        <ToggleSwitch
                          checked={card.openBankingEnabled}
                          onChange={(checked) => updateCreditCard(card.id, 'openBankingEnabled', checked)}
                        />
                      </div>

                      {card.openBankingEnabled && (
                        <div className="space-y-4 p-4 bg-gray-600 rounded-lg">
                          <h4 className="font-medium text-sm">TrueLayer Open Banking Configuration</h4>
                          
                          <div>
                            <Label className="flex items-center justify-between">
                              Client ID
                              {getStatusBadge(true, !!card.truelayerClientId)}
                            </Label>
                            <TextInput
                              value={card.truelayerClientId}
                              onChange={(e) => updateCreditCard(card.id, 'truelayerClientId', e.target.value)}
                              placeholder="your-truelayer-client-id"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>

                          <div>
                            <Label className="flex items-center justify-between">
                              Client Secret
                              {getStatusBadge(true, !!card.truelayerClientSecret)}
                            </Label>
                            <div className="relative">
                              <TextInput
                                type={showOAuthSecret ? "text" : "password"}
                                value={card.truelayerClientSecret}
                                onChange={(e) => updateCreditCard(card.id, 'truelayerClientSecret', e.target.value)}
                                placeholder="your-truelayer-client-secret"
                                className="bg-gray-500 border-gray-400 text-white pr-10"
                              />
                              <button
                                type="button"
                                className="absolute right-2 top-1/2 transform -translate-y-1/2"
                                onClick={() => setShowOAuthSecret(!showOAuthSecret)}
                              >
                                {showOAuthSecret ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                              </button>
                            </div>
                          </div>

                          <div>
                            <Label className="flex items-center justify-between">
                              Redirect URI
                              {getStatusBadge(true, !!card.truelayerRedirectUri)}
                            </Label>
                            <TextInput
                              value={card.truelayerRedirectUri}
                              onChange={(e) => updateCreditCard(card.id, 'truelayerRedirectUri', e.target.value)}
                              placeholder="https://yourapp.com/oauth/callback"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>

                          <div>
                            <Label className="flex items-center justify-between">
                              Environment
                              {getStatusBadge(false, true)}
                            </Label>
                            <Select 
                              value={card.truelayerEnvironment} 
                              onChange={(e) => updateCreditCard(card.id, 'truelayerEnvironment', e.target.value)}
                            >
                              <option value="sandbox">Sandbox (Testing)</option>
                              <option value="live">Live (Production)</option>
                            </Select>
                          </div>
                        </div>
                      )}
                    </div>
                  </Card>
                ))}
              </div>
            </TabItem>

            <TabItem title={
              <div className="flex items-center gap-2">
                <MessageSquare className="w-4 h-4" />
                Communication
                <Badge color="gray" className="ml-auto">{communicationChannels.length}</Badge>
              </div>
            }>
              <div className="p-6 space-y-4">
                <h2 className="text-lg font-semibold">Communication Channels</h2>
                
                <Card className="bg-gray-700 border-gray-600">
                  <div className="flex items-center gap-2 pb-3">
                    <Mail className="w-4 h-4" />
                    <h3 className="font-semibold">Email Configuration</h3>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label className="flex items-center justify-between">
                        Email Provider
                        {getStatusBadge(true, !!emailConfig.provider)}
                      </Label>
                      <Select 
                        value={emailConfig.provider}
                        onChange={(e) => updateEmailConfig('provider', e.target.value)}
                      >
                        <option value="">Select provider</option>
                        <option value="gmail">Gmail</option>
                        <option value="outlook">Outlook/Microsoft 365</option>
                        <option value="yahoo">Yahoo Mail</option>
                        <option value="custom">Custom IMAP/SMTP</option>
                      </Select>
                    </div>
                    
                    <div>
                      <Label className="flex items-center justify-between">
                        Email Address
                        {getStatusBadge(true, !!emailConfig.emailAddress)}
                      </Label>
                      <TextInput
                        value={emailConfig.emailAddress}
                        onChange={(e) => updateEmailConfig('emailAddress', e.target.value)}
                        placeholder="your@email.com"
                        className="bg-gray-600 border-gray-500 text-white"
                      />
                    </div>

                    <div>
                      <Label className="flex items-center justify-between">
                        OAuth2 Enabled
                        {getStatusBadge(false, true)}
                      </Label>
                      <ToggleSwitch
                        checked={emailConfig.oauth2Enabled}
                        onChange={(checked) => updateEmailConfig('oauth2Enabled', checked)}
                      />
                    </div>

                    {emailConfig.oauth2Enabled ? (
                      <div className="space-y-4 p-4 bg-gray-600 rounded-lg">
                        <h4 className="font-medium text-sm">OAuth2 Configuration</h4>
                        
                        <div>
                          <Label className="flex items-center justify-between">
                            OAuth2 Client ID
                            {getStatusBadge(true, !!emailConfig.oauth2ClientId)}
                          </Label>
                          <TextInput
                            value={emailConfig.oauth2ClientId}
                            onChange={(e) => updateEmailConfig('oauth2ClientId', e.target.value)}
                            placeholder="your-oauth2-client-id"
                            className="bg-gray-500 border-gray-400 text-white"
                          />
                        </div>

                        <div>
                          <Label className="flex items-center justify-between">
                            OAuth2 Client Secret
                            {getStatusBadge(true, !!emailConfig.oauth2ClientSecret)}
                          </Label>
                          <div className="relative">
                            <TextInput
                              type={showOAuthSecret ? "text" : "password"}
                              value={emailConfig.oauth2ClientSecret}
                              onChange={(e) => updateEmailConfig('oauth2ClientSecret', e.target.value)}
                              placeholder="your-oauth2-client-secret"
                              className="bg-gray-500 border-gray-400 text-white pr-10"
                            />
                            <button
                              type="button"
                              className="absolute right-2 top-1/2 transform -translate-y-1/2"
                              onClick={() => setShowOAuthSecret(!showOAuthSecret)}
                            >
                              {showOAuthSecret ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                            </button>
                          </div>
                        </div>

                        <div>
                          <Label className="flex items-center justify-between">
                            Refresh Token
                            {getStatusBadge(false, !!emailConfig.oauth2RefreshToken)}
                          </Label>
                          <Textarea
                            value={emailConfig.oauth2RefreshToken}
                            onChange={(e) => updateEmailConfig('oauth2RefreshToken', e.target.value)}
                            placeholder="OAuth2 refresh token for automatic token renewal"
                            className="bg-gray-500 border-gray-400 text-white"
                          />
                        </div>

                        <div>
                          <Label className="flex items-center justify-between">
                            Access Token
                            {getStatusBadge(false, !!emailConfig.oauth2AccessToken)}
                          </Label>
                          <Textarea
                            value={emailConfig.oauth2AccessToken}
                            onChange={(e) => updateEmailConfig('oauth2AccessToken', e.target.value)}
                            placeholder="OAuth2 access token for email API access"
                            className="bg-gray-500 border-gray-400 text-white"
                          />
                        </div>
                      </div>
                    ) : (
                      <div>
                        <Label className="flex items-center justify-between">
                          Password
                          {getStatusBadge(true, !!emailConfig.password)}
                        </Label>
                        <div className="relative">
                          <TextInput
                            type={showPassword ? "text" : "password"}
                            value={emailConfig.password}
                            onChange={(e) => updateEmailConfig('password', e.target.value)}
                            placeholder="Email account password"
                            className="bg-gray-600 border-gray-500 text-white pr-10"
                          />
                          <button
                            type="button"
                            className="absolute right-2 top-1/2 transform -translate-y-1/2"
                            onClick={() => setShowPassword(!showPassword)}
                          >
                            {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                          </button>
                        </div>
                      </div>
                    )}

                    <div>
                      <Label className="flex items-center justify-between">
                        Auto Sync Enabled
                        {getStatusBadge(false, true)}
                      </Label>
                      <ToggleSwitch
                        checked={emailConfig.autoSyncEnabled}
                        onChange={(checked) => updateEmailConfig('autoSyncEnabled', checked)}
                      />
                    </div>

                    {emailConfig.autoSyncEnabled && (
                      <div>
                        <Label className="flex items-center justify-between">
                          Sync Frequency
                          {getStatusBadge(false, true)}
                        </Label>
                        <Select 
                          value={emailConfig.syncFrequency}
                          onChange={(e) => updateEmailConfig('syncFrequency', e.target.value)}
                        >
                          <option value="1min">Every minute</option>
                          <option value="5min">Every 5 minutes</option>
                          <option value="15min">Every 15 minutes</option>
                          <option value="30min">Every 30 minutes</option>
                          <option value="1hour">Every hour</option>
                        </Select>
                      </div>
                    )}

                    {emailConfig.provider === 'custom' && (
                      <div className="space-y-4 p-4 bg-gray-600 rounded-lg">
                        <h4 className="font-medium text-sm">Custom Server Configuration</h4>
                        
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <Label className="flex items-center justify-between">
                              IMAP Server
                              {getStatusBadge(true, !!emailConfig.imapServer)}
                            </Label>
                            <TextInput
                              value={emailConfig.imapServer}
                              onChange={(e) => updateEmailConfig('imapServer', e.target.value)}
                              placeholder="imap.example.com"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>
                          <div>
                            <Label className="flex items-center justify-between">
                              IMAP Port
                              {getStatusBadge(false, true)}
                            </Label>
                            <TextInput
                              value={emailConfig.imapPort}
                              onChange={(e) => updateEmailConfig('imapPort', e.target.value)}
                              placeholder="993"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <Label className="flex items-center justify-between">
                              SMTP Server
                              {getStatusBadge(true, !!emailConfig.smtpServer)}
                            </Label>
                            <TextInput
                              value={emailConfig.smtpServer}
                              onChange={(e) => updateEmailConfig('smtpServer', e.target.value)}
                              placeholder="smtp.example.com"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>
                          <div>
                            <Label className="flex items-center justify-between">
                              SMTP Port
                              {getStatusBadge(false, true)}
                            </Label>
                            <TextInput
                              value={emailConfig.smtpPort}
                              onChange={(e) => updateEmailConfig('smtpPort', e.target.value)}
                              placeholder="587"
                              className="bg-gray-500 border-gray-400 text-white"
                            />
                          </div>
                        </div>

                        <div className="flex gap-4">
                          <div className="flex items-center gap-2">
                            <ToggleSwitch
                              checked={emailConfig.useSSL}
                              onChange={(checked) => updateEmailConfig('useSSL', checked)}
                            />
                            <Label>Use SSL</Label>
                          </div>
                          <div className="flex items-center gap-2">
                            <ToggleSwitch
                              checked={emailConfig.useTLS}
                              onChange={(checked) => updateEmailConfig('useTLS', checked)}
                            />
                            <Label>Use TLS</Label>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                </Card>

                <Card className="bg-gray-700 border-gray-600">
                  <div className="flex items-center gap-2 pb-3">
                    <MessageSquare className="w-4 h-4" />
                    <h3 className="font-semibold">SMS Configuration</h3>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label className="flex items-center justify-between">
                        Mobile Number
                        {getStatusBadge(false, false)}
                      </Label>
                      <TextInput
                        placeholder="+44 123 456 7890"
                        className="bg-gray-600 border-gray-500 text-white"
                      />
                    </div>
                  </div>
                </Card>

                <Card className="bg-gray-700 border-gray-600">
                  <div className="flex items-center gap-2 pb-3">
                    <Mic className="w-4 h-4" />
                    <h3 className="font-semibold">Voice Configuration</h3>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label className="flex items-center justify-between">
                        Voice Assistant
                        {getStatusBadge(false, false)}
                      </Label>
                      <Select>
                        <option value="">Select assistant</option>
                        <option value="alexa">Amazon Alexa</option>
                        <option value="google">Google Assistant</option>
                        <option value="siri">Apple Siri</option>
                        <option value="cortana">Microsoft Cortana</option>
                      </Select>
                    </div>
                  </div>
                </Card>

                <Card className="bg-gray-700 border-gray-600">
                  <div className="flex items-center gap-2 pb-3">
                    <Lock className="w-4 h-4" />
                    <h3 className="font-semibold">Secure Messenger</h3>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label className="flex items-center justify-between">
                        Messenger Type
                        {getStatusBadge(false, false)}
                      </Label>
                      <Select>
                        <option value="">Select messenger</option>
                        <option value="telegram">Telegram</option>
                        <option value="signal">Signal</option>
                        <option value="whatsapp">WhatsApp Business</option>
                      </Select>
                    </div>
                  </div>
                </Card>
              </div>
            </TabItem>

            <TabItem title={
              <div className="flex items-center gap-2">
                <Shield className="w-4 h-4" />
                Security
                <Badge color="gray" className="ml-auto">3</Badge>
              </div>
            }>
              <div className="p-6 space-y-4">
                <h2 className="text-lg font-semibold">Security Settings</h2>
                
                <Card className="bg-gray-700 border-gray-600">
                  <div className="flex items-center gap-2 pb-3">
                    <Shield className="w-4 h-4" />
                    <h3 className="font-semibold">Authentication</h3>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <Label className="flex items-center justify-between">
                        Two-Factor Authentication
                        {getStatusBadge(false, false)}
                      </Label>
                      <Select>
                        <option value="">Select 2FA method</option>
                        <option value="none">Disabled</option>
                        <option value="sms">SMS</option>
                        <option value="authenticator">Authenticator App</option>
                        <option value="hardware">Hardware Token</option>
                      </Select>
                    </div>
                  </div>
                </Card>
              </div>
            </TabItem>
          </Tabs>
        </div>
      </div>

      {/* Floating Upload Documents Modal */}
      <Modal show={isUploadOpen} onClose={() => setIsUploadOpen(false)} size="2xl" className="z-50">
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Upload Documents</h3>
            <Button color="gray" size="sm" onClick={() => setIsUploadOpen(false)}>
              <X className="w-4 h-4" />
            </Button>
          </div>
          <div className="space-y-4">
            <div className="border-2 border-dashed border-gray-600 rounded-lg p-8 text-center">
              <Upload className="w-12 h-12 mx-auto text-gray-400 mb-4" />
              <p className="text-gray-400 mb-2">Upload or drag and drop documents</p>
              <p className="text-sm text-gray-500">Supports: PDF, JPG, PNG, TTL</p>
              <Button className="mt-4">
                <Upload className="w-4 h-4 mr-2" />
                Browse Files
              </Button>
            </div>
            <div className="space-y-2">
              <h4 className="font-medium">Supported Documents:</h4>
              <ul className="text-sm text-gray-400 space-y-1">
                <li>• Driving License</li>
                <li>• Passport</li>
                <li>• Bank Statements</li>
                <li>• Credit Card Photos</li>
                <li>• Insurance Quotes</li>
                <li>• Configuration Screenshots</li>
              </ul>
            </div>
            <div className="bg-gray-700 rounded-lg p-4">
              <h4 className="font-medium mb-2">Upload History</h4>
              <div className="text-sm text-gray-400">
                <div className="flex items-center justify-between py-2 border-b border-gray-600">
                  <span>driving_license.pdf</span>
                  <div className="flex gap-2">
                    <Button size="xs" color="gray">
                      <Download className="w-3 h-3" />
                    </Button>
                    <Button size="xs" color="failure">
                      <Trash2 className="w-3 h-3" />
                    </Button>
                  </div>
                </div>
                <div className="flex items-center justify-between py-2">
                  <span>bank_statement.pdf</span>
                  <div className="flex gap-2">
                    <Button size="xs" color="gray">
                      <Download className="w-3 h-3" />
                    </Button>
                    <Button size="xs" color="failure">
                      <Trash2 className="w-3 h-3" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex justify-end gap-2 mt-6">
            <Button color="gray" onClick={() => setIsUploadOpen(false)}>Cancel</Button>
            <Button>Upload</Button>
          </div>
        </div>
      </Modal>

      {/* Floating Help Modal */}
      <Modal show={isHelpOpen} onClose={() => setIsHelpOpen(false)} size="2xl" className="z-50">
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">AI Assistant - Help</h3>
            <Button color="gray" size="sm" onClick={() => setIsHelpOpen(false)}>
              <X className="w-4 h-4" />
            </Button>
          </div>
          <div className="space-y-4">
            <div className="bg-gray-700 rounded-lg p-4">
              <p className="text-sm">
                Welcome! I'm here to help you configure your settings. I can assist with:
              </p>
              <ul className="text-sm text-gray-300 mt-2 space-y-1">
                <li>• Setting up bank accounts and Open Banking</li>
                <li>• Adding credit cards securely</li>
                <li>• Configuring communication channels</li>
                <li>• Security settings and 2FA</li>
                <li>• Document upload and field extraction</li>
              </ul>
            </div>
            <div className="bg-gray-700 rounded-lg p-4">
              <h4 className="font-medium mb-2">Recent Conversations</h4>
              <div className="text-sm text-gray-400 space-y-2">
                <div className="p-2 bg-gray-600 rounded">
                  <div className="font-medium text-white">You:</div>
                  <div>How do I set up OAuth2 for Gmail?</div>
                </div>
                <div className="p-2 bg-gray-600 rounded">
                  <div className="font-medium text-white">AI Assistant:</div>
                  <div>To set up OAuth2 for Gmail, you need to: 1. Go to Google Cloud Console 2. Create a project 3. Enable Gmail API 4. Create OAuth2 credentials 5. Add the client ID and secret to your email configuration.</div>
                </div>
              </div>
            </div>
            <div className="flex gap-2">
              <TextInput placeholder="Ask me anything..." className="bg-gray-600 border-gray-500 text-white flex-1" />
              <Button>Send</Button>
            </div>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default App;
