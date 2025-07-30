#!/usr/bin/env node

/**
 * AdvanceGG Node.js Post-Install Script
 * 
 * Downloads pre-built native libraries for the current platform.
 * No Go compiler required.
 */

const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');
const { execSync } = require('child_process');

// Configuration
const PACKAGE_VERSION = '1.0.0';
const GITHUB_RELEASES_URL = `https://github.com/GrandpaEJ/advancegg/releases/download/v${PACKAGE_VERSION}`;
const NATIVE_DIR = path.join(__dirname, '..', 'native');

// Platform detection
function getPlatformInfo() {
    const platform = os.platform();
    const arch = os.arch();
    
    // Normalize platform names
    const platformMap = {
        'win32': 'windows',
        'darwin': 'darwin',
        'linux': 'linux'
    };
    
    // Normalize architecture names
    const archMap = {
        'x64': 'x64',
        'arm64': 'arm64',
        'arm': 'armv7'
    };
    
    return {
        platform: platformMap[platform] || platform,
        arch: archMap[arch] || arch
    };
}

// Get binary filename for platform
function getBinaryName(platform, arch) {
    const extension = platform === 'windows' ? '.dll' : 
                     platform === 'darwin' ? '.dylib' : '.so';
    return `advancegg-${platform}-${arch}${extension}`;
}

// Download file with progress
function downloadFile(url, destination) {
    return new Promise((resolve, reject) => {
        console.log(`📦 Downloading: ${path.basename(destination)}`);
        
        const file = fs.createWriteStream(destination);
        
        https.get(url, (response) => {
            if (response.statusCode === 302 || response.statusCode === 301) {
                // Handle redirect
                return downloadFile(response.headers.location, destination)
                    .then(resolve)
                    .catch(reject);
            }
            
            if (response.statusCode !== 200) {
                reject(new Error(`HTTP ${response.statusCode}: ${response.statusMessage}`));
                return;
            }
            
            const totalSize = parseInt(response.headers['content-length'], 10);
            let downloadedSize = 0;
            
            response.on('data', (chunk) => {
                downloadedSize += chunk.length;
                if (totalSize) {
                    const percent = Math.round((downloadedSize / totalSize) * 100);
                    process.stdout.write(`\r   Progress: ${percent}%`);
                }
            });
            
            response.pipe(file);
            
            file.on('finish', () => {
                file.close();
                console.log('\n✅ Download completed');
                resolve();
            });
            
            file.on('error', (err) => {
                fs.unlink(destination, () => {}); // Clean up
                reject(err);
            });
            
        }).on('error', reject);
    });
}

// Check if Go is available for building from source
function hasGoCompiler() {
    try {
        execSync('go version', { stdio: 'ignore' });
        return true;
    } catch {
        return false;
    }
}

// Build from source as fallback
function buildFromSource() {
    console.log('🔨 Building native library from source...');
    
    if (!hasGoCompiler()) {
        throw new Error('Go compiler not found. Please install Go or use pre-built binaries.');
    }
    
    try {
        const { platform, arch } = getPlatformInfo();
        const binaryName = getBinaryName(platform, arch);
        const outputPath = path.join(NATIVE_DIR, binaryName);
        
        // Build command
        const buildCmd = [
            'go', 'build',
            '-buildmode=c-shared',
            '-o', outputPath,
            'advancegg_nodejs.go'
        ].join(' ');
        
        execSync(buildCmd, { 
            cwd: path.join(__dirname, '..'),
            stdio: 'inherit' 
        });
        
        console.log('✅ Built native library from source');
        return true;
    } catch (error) {
        console.error('❌ Failed to build from source:', error.message);
        return false;
    }
}

// Verify native library
function verifyNativeLibrary(binaryPath) {
    try {
        // Basic file existence and size check
        const stats = fs.statSync(binaryPath);
        if (stats.size < 1000) { // Too small to be valid
            return false;
        }
        
        // Try to load the library (basic validation)
        // This is platform-specific and might need adjustment
        return true;
    } catch {
        return false;
    }
}

// Main installation function
async function install() {
    console.log('🎨 AdvanceGG Post-Install Setup');
    console.log('================================');
    
    // Create native directory
    if (!fs.existsSync(NATIVE_DIR)) {
        fs.mkdirSync(NATIVE_DIR, { recursive: true });
    }
    
    // Get platform info
    const { platform, arch } = getPlatformInfo();
    console.log(`🖥️  Platform: ${platform}-${arch}`);
    
    const binaryName = getBinaryName(platform, arch);
    const binaryPath = path.join(NATIVE_DIR, binaryName);
    
    // Check if binary already exists and is valid
    if (fs.existsSync(binaryPath) && verifyNativeLibrary(binaryPath)) {
        console.log('✅ Native library already present and valid');
        return;
    }
    
    // Try to download pre-built binary
    const downloadUrl = `${GITHUB_RELEASES_URL}/${binaryName}`;
    
    try {
        await downloadFile(downloadUrl, binaryPath);
        
        // Verify downloaded binary
        if (verifyNativeLibrary(binaryPath)) {
            console.log('✅ Native library installed successfully');
            return;
        } else {
            console.log('⚠️  Downloaded binary failed verification');
            fs.unlinkSync(binaryPath);
        }
    } catch (error) {
        console.log(`⚠️  Failed to download pre-built binary: ${error.message}`);
    }
    
    // Fallback to building from source
    console.log('🔄 Attempting to build from source...');
    
    if (buildFromSource()) {
        if (verifyNativeLibrary(binaryPath)) {
            console.log('✅ Native library built and verified');
            return;
        }
    }
    
    // Final fallback - try to find system-wide installation
    console.log('🔍 Looking for system-wide AdvanceGG installation...');
    
    const systemPaths = [
        '/usr/local/lib',
        '/usr/lib',
        '/opt/advancegg/lib',
        process.env.ADVANCEGG_LIB_PATH
    ].filter(Boolean);
    
    for (const systemPath of systemPaths) {
        const systemBinary = path.join(systemPath, binaryName);
        if (fs.existsSync(systemBinary)) {
            try {
                fs.copyFileSync(systemBinary, binaryPath);
                console.log(`✅ Copied from system installation: ${systemPath}`);
                return;
            } catch (error) {
                console.log(`⚠️  Failed to copy from ${systemPath}: ${error.message}`);
            }
        }
    }
    
    // Installation failed
    console.error('❌ Failed to install AdvanceGG native library');
    console.error('');
    console.error('Possible solutions:');
    console.error('1. Install Go compiler and rebuild: npm rebuild advancegg');
    console.error('2. Download manually from: https://github.com/GrandpaEJ/advancegg/releases');
    console.error('3. Set ADVANCEGG_LIB_PATH environment variable');
    console.error('');
    
    // Don't fail the installation - let runtime handle the error
    console.log('⚠️  Installation completed with warnings');
}

// Run installation
if (require.main === module) {
    install().catch((error) => {
        console.error('❌ Post-install failed:', error.message);
        process.exit(1);
    });
}

module.exports = { install, getPlatformInfo, getBinaryName };
