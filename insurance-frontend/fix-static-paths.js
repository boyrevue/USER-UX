#!/usr/bin/env node

/**
 * ============================================================
 * STATIC PATH FIXER FOR CLIENT-UX
 * ============================================================
 * Permanent solution to fix the static/static nested directory issue
 * This script runs automatically after React build to ensure
 * proper static file paths and prevent MIME type issues
 * ============================================================
 */

const fs = require('fs');
const path = require('path');

console.log('🔧 Fixing static paths for CLIENT-UX...');

const buildDir = path.join(__dirname, 'build');
const indexPath = path.join(buildDir, 'index.html');
const assetManifestPath = path.join(buildDir, 'asset-manifest.json');

/**
 * Fix index.html to use correct static paths
 */
function fixIndexHtml() {
    if (!fs.existsSync(indexPath)) {
        console.log('❌ index.html not found');
        return false;
    }

    let indexContent = fs.readFileSync(indexPath, 'utf8');
    
    // Fix static paths in HTML
    const originalContent = indexContent;
    
    // Replace /static/ with /static/ (ensure no double static)
    indexContent = indexContent.replace(/\/static\/static\//g, '/static/');
    
    // Ensure all static references are correct
    indexContent = indexContent.replace(/href="\/static\//g, 'href="/static/');
    indexContent = indexContent.replace(/src="\/static\//g, 'src="/static/');
    
    if (indexContent !== originalContent) {
        fs.writeFileSync(indexPath, indexContent);
        console.log('✅ Fixed index.html static paths');
        return true;
    } else {
        console.log('✅ index.html paths already correct');
        return true;
    }
}

/**
 * Fix asset-manifest.json paths
 */
function fixAssetManifest() {
    if (!fs.existsSync(assetManifestPath)) {
        console.log('❌ asset-manifest.json not found');
        return false;
    }

    try {
        const manifest = JSON.parse(fs.readFileSync(assetManifestPath, 'utf8'));
        let modified = false;

        // Fix all file paths in manifest
        for (const key in manifest.files) {
            if (manifest.files[key].includes('/static/static/')) {
                manifest.files[key] = manifest.files[key].replace('/static/static/', '/static/');
                modified = true;
            }
        }

        // Fix entrypoints
        if (manifest.entrypoints) {
            manifest.entrypoints = manifest.entrypoints.map(entry => {
                if (entry.includes('/static/static/')) {
                    return entry.replace('/static/static/', '/static/');
                }
                return entry;
            });
            modified = true;
        }

        if (modified) {
            fs.writeFileSync(assetManifestPath, JSON.stringify(manifest, null, 2));
            console.log('✅ Fixed asset-manifest.json paths');
        } else {
            console.log('✅ asset-manifest.json paths already correct');
        }

        return true;
    } catch (error) {
        console.log('❌ Error fixing asset-manifest.json:', error.message);
        return false;
    }
}

/**
 * Verify static directory structure
 */
function verifyStaticStructure() {
    const staticDir = path.join(buildDir, 'static');
    const nestedStaticDir = path.join(staticDir, 'static');

    if (fs.existsSync(nestedStaticDir)) {
        console.log('⚠️  Found nested static/static directory - this should not happen');
        console.log('   This indicates a build configuration issue');
        return false;
    }

    // Check for required files
    const cssDir = path.join(staticDir, 'css');
    const jsDir = path.join(staticDir, 'js');

    if (!fs.existsSync(cssDir)) {
        console.log('❌ CSS directory not found in static/');
        return false;
    }

    if (!fs.existsSync(jsDir)) {
        console.log('❌ JS directory not found in static/');
        return false;
    }

    const cssFiles = fs.readdirSync(cssDir).filter(f => f.endsWith('.css'));
    const jsFiles = fs.readdirSync(jsDir).filter(f => f.endsWith('.js') && !f.includes('.map'));

    if (cssFiles.length === 0) {
        console.log('❌ No CSS files found');
        return false;
    }

    if (jsFiles.length === 0) {
        console.log('❌ No JS files found');
        return false;
    }

    console.log(`✅ Static structure verified: ${cssFiles.length} CSS, ${jsFiles.length} JS files`);
    return true;
}

/**
 * Create deployment summary
 */
function createDeploymentSummary() {
    const staticDir = path.join(buildDir, 'static');
    const summary = {
        timestamp: new Date().toISOString(),
        build_dir: buildDir,
        static_files: {},
        status: 'ready'
    };

    try {
        // Get CSS files
        const cssDir = path.join(staticDir, 'css');
        if (fs.existsSync(cssDir)) {
            summary.static_files.css = fs.readdirSync(cssDir)
                .filter(f => f.endsWith('.css'))
                .map(f => `/static/css/${f}`);
        }

        // Get JS files
        const jsDir = path.join(staticDir, 'js');
        if (fs.existsSync(jsDir)) {
            summary.static_files.js = fs.readdirSync(jsDir)
                .filter(f => f.endsWith('.js') && !f.includes('.map'))
                .map(f => `/static/js/${f}`);
        }

        const summaryPath = path.join(buildDir, 'build-summary.json');
        fs.writeFileSync(summaryPath, JSON.stringify(summary, null, 2));
        console.log('✅ Created build summary');

    } catch (error) {
        console.log('⚠️  Could not create build summary:', error.message);
    }
}

/**
 * Main execution
 */
function main() {
    console.log('Starting static path fixes...');
    
    let success = true;
    
    success &= fixIndexHtml();
    success &= fixAssetManifest();
    success &= verifyStaticStructure();
    
    createDeploymentSummary();
    
    if (success) {
        console.log('🎉 All static path fixes completed successfully!');
        console.log('📁 Build directory:', buildDir);
        console.log('🔗 Static files ready for deployment');
        process.exit(0);
    } else {
        console.log('❌ Some fixes failed - check output above');
        process.exit(1);
    }
}

// Run the fixer
main();
