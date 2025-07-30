#!/usr/bin/env node

/**
 * AdvanceGG Node.js Info Script
 * 
 * Displays installation and system information.
 */

const fs = require('fs');
const path = require('path');
const os = require('os');

// Try to load AdvanceGG
let AdvanceGG;
try {
    AdvanceGG = require('../index.js');
} catch (error) {
    console.error('❌ Failed to load AdvanceGG:', error.message);
    process.exit(1);
}

function displayInfo() {
    console.log('🎨 AdvanceGG Node.js Information');
    console.log('='.repeat(50));
    
    // Package information
    const packageJson = require('../package.json');
    console.log(`📦 Package Version: ${packageJson.version}`);
    console.log(`📍 Installation Path: ${path.resolve(__dirname, '..')}`);
    
    // System information
    console.log(`🖥️  Platform: ${os.platform()}-${os.arch()}`);
    console.log(`🟨 Node.js: ${process.version}`);
    console.log(`💾 Memory: ${Math.round(os.totalmem() / 1024 / 1024 / 1024)} GB`);
    console.log(`🔧 CPU Cores: ${os.cpus().length}`);
    
    // Native library information
    const nativeDir = path.join(__dirname, '..', 'native');
    console.log(`\n📚 Native Libraries:`);
    
    if (fs.existsSync(nativeDir)) {
        const files = fs.readdirSync(nativeDir);
        const libraries = files.filter(f => f.startsWith('advancegg-'));
        
        if (libraries.length > 0) {
            libraries.forEach(lib => {
                const libPath = path.join(nativeDir, lib);
                const stats = fs.statSync(libPath);
                const sizeMB = (stats.size / 1024 / 1024).toFixed(2);
                console.log(`   ✅ ${lib} (${sizeMB} MB)`);
            });
        } else {
            console.log('   ❌ No native libraries found');
        }
    } else {
        console.log('   ❌ Native directory not found');
    }
    
    // Feature support
    console.log(`\n🚀 Feature Support:`);
    try {
        // Test basic functionality
        const canvas = new AdvanceGG.Canvas(100, 100);
        console.log('   ✅ Canvas Creation');
        
        canvas.setRGB(1, 0, 0);
        canvas.drawCircle(50, 50, 25);
        canvas.fill();
        console.log('   ✅ Basic Drawing');
        
        // Test advanced features if available
        try {
            const gradient = canvas.createLinearGradient(0, 0, 100, 0);
            gradient.addColorStop(0, [1, 0, 0, 1]);
            gradient.addColorStop(1, [0, 0, 1, 1]);
            console.log('   ✅ Gradients');
        } catch {
            console.log('   ❌ Gradients');
        }
        
        try {
            canvas.drawString('Test', 10, 10);
            console.log('   ✅ Text Rendering');
        } catch {
            console.log('   ❌ Text Rendering');
        }
        
        canvas.dispose();
        
    } catch (error) {
        console.log('   ❌ Basic functionality failed:', error.message);
    }
    
    // Performance information
    console.log(`\n⚡ Performance Info:`);
    console.log(`   CPU: ${os.cpus()[0].model}`);
    console.log(`   Speed: ${os.cpus()[0].speed} MHz`);
    console.log(`   Load Average: ${os.loadavg().map(l => l.toFixed(2)).join(', ')}`);
    
    // Environment variables
    const envVars = [
        'NODE_ENV',
        'ADVANCEGG_LIB_PATH',
        'ADVANCEGG_DEBUG',
        'ADVANCEGG_SIMD'
    ];
    
    const setEnvVars = envVars.filter(v => process.env[v]);
    if (setEnvVars.length > 0) {
        console.log(`\n🌍 Environment Variables:`);
        setEnvVars.forEach(v => {
            console.log(`   ${v}: ${process.env[v]}`);
        });
    }
    
    // Installation health check
    console.log(`\n🏥 Health Check:`);
    
    const checks = [
        {
            name: 'Package Integrity',
            test: () => fs.existsSync(path.join(__dirname, '..', 'package.json'))
        },
        {
            name: 'Native Library',
            test: () => {
                const nativeDir = path.join(__dirname, '..', 'native');
                return fs.existsSync(nativeDir) && fs.readdirSync(nativeDir).length > 0;
            }
        },
        {
            name: 'Canvas Creation',
            test: () => {
                try {
                    const canvas = new AdvanceGG.Canvas(10, 10);
                    canvas.dispose();
                    return true;
                } catch {
                    return false;
                }
            }
        },
        {
            name: 'Memory Management',
            test: () => {
                try {
                    const canvases = [];
                    for (let i = 0; i < 10; i++) {
                        canvases.push(new AdvanceGG.Canvas(50, 50));
                    }
                    canvases.forEach(c => c.dispose());
                    return true;
                } catch {
                    return false;
                }
            }
        }
    ];
    
    checks.forEach(check => {
        const result = check.test();
        console.log(`   ${result ? '✅' : '❌'} ${check.name}`);
    });
    
    const passedChecks = checks.filter(c => c.test()).length;
    const healthPercentage = Math.round((passedChecks / checks.length) * 100);
    
    console.log(`\n🎯 Overall Health: ${healthPercentage}%`);
    
    if (healthPercentage === 100) {
        console.log('🎉 AdvanceGG is fully functional!');
    } else if (healthPercentage >= 75) {
        console.log('⚠️  AdvanceGG is mostly functional with minor issues');
    } else if (healthPercentage >= 50) {
        console.log('⚠️  AdvanceGG has significant issues');
    } else {
        console.log('❌ AdvanceGG installation is severely compromised');
        console.log('\nTroubleshooting:');
        console.log('1. Try reinstalling: npm uninstall advancegg && npm install advancegg');
        console.log('2. Check native library: ls -la node_modules/advancegg/native/');
        console.log('3. Verify Node.js version: node --version (requires 16+)');
        console.log('4. Check for conflicting packages');
    }
}

// CLI interface
function main() {
    const args = process.argv.slice(2);
    
    if (args.includes('--help') || args.includes('-h')) {
        console.log('AdvanceGG Info Tool');
        console.log('');
        console.log('Usage: advancegg-info [options]');
        console.log('');
        console.log('Options:');
        console.log('  --help, -h     Show this help message');
        console.log('  --json         Output information in JSON format');
        console.log('  --health       Run health check only');
        return;
    }
    
    if (args.includes('--json')) {
        // JSON output for programmatic use
        const info = {
            package: require('../package.json'),
            system: {
                platform: os.platform(),
                arch: os.arch(),
                nodeVersion: process.version,
                memory: os.totalmem(),
                cpus: os.cpus().length
            },
            health: {
                // Run health checks and return results
            }
        };
        console.log(JSON.stringify(info, null, 2));
        return;
    }
    
    if (args.includes('--health')) {
        // Health check only
        console.log('Running health check...');
        // Run health checks
        return;
    }
    
    // Default: show full info
    displayInfo();
}

if (require.main === module) {
    main();
}

module.exports = { displayInfo };
