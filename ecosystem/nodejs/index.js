/**
 * AdvanceGG Node.js WebAssembly Wrapper
 * 
 * High-performance 2D graphics library for Node.js using Go WebAssembly backend.
 */

const fs = require('fs');
const path = require('path');

// Global WASM instance and Go runtime
let wasmInstance = null;
let go = null;
let isInitialized = false;

// Context storage
const contexts = new Map();
let nextContextId = 1;

/**
 * Initialize the WebAssembly module
 */
async function init() {
    if (isInitialized) return;
    
    try {
        // Load Go WASM runtime
        const wasmExecPath = path.join(__dirname, 'wasm_exec.js');
        if (!fs.existsSync(wasmExecPath)) {
            throw new Error('wasm_exec.js not found. Please ensure Go WebAssembly support files are present.');
        }
        
        // Import Go WASM runtime
        require('./wasm_exec.js');
        go = new Go();
        
        // Load WASM module
        const wasmPath = path.join(__dirname, 'advancegg.wasm');
        if (!fs.existsSync(wasmPath)) {
            throw new Error('advancegg.wasm not found. Please build the WebAssembly module first.');
        }
        
        const wasmBuffer = fs.readFileSync(wasmPath);
        const wasmModule = await WebAssembly.instantiate(wasmBuffer, go.importObject);
        
        wasmInstance = wasmModule.instance;
        
        // Start Go runtime
        go.run(wasmInstance);
        
        isInitialized = true;
        console.log('AdvanceGG WebAssembly module initialized successfully');
    } catch (error) {
        throw new Error(`Failed to initialize AdvanceGG WASM: ${error.message}`);
    }
}

/**
 * Ensure WASM is initialized
 */
function ensureInitialized() {
    if (!isInitialized) {
        throw new Error('AdvanceGG not initialized. Call await AdvanceGG.init() first.');
    }
}

/**
 * Canvas class for 2D graphics operations
 */
class Canvas {
    constructor(width, height) {
        ensureInitialized();
        
        this.width = width;
        this.height = height;
        this._id = nextContextId++;
        
        // Call Go function to create context
        if (global.createContext) {
            const contextId = global.createContext(width, height);
            contexts.set(this._id, contextId);
        } else {
            throw new Error('createContext function not available from WASM module');
        }
    }
    
    /**
     * Clean up the canvas
     */
    dispose() {
        const contextId = contexts.get(this._id);
        if (contextId && global.destroyContext) {
            global.destroyContext(contextId);
            contexts.delete(this._id);
        }
    }
    
    /**
     * Set RGB color
     */
    setRGB(r, g, b) {
        const contextId = contexts.get(this._id);
        if (contextId && global.setRGB) {
            global.setRGB(contextId, r, g, b);
        }
    }
    
    /**
     * Set RGBA color
     */
    setRGBA(r, g, b, a) {
        const contextId = contexts.get(this._id);
        if (contextId && global.setRGBA) {
            global.setRGBA(contextId, r, g, b, a);
        }
    }
    
    /**
     * Set hex color
     */
    setHexColor(hexColor) {
        const contextId = contexts.get(this._id);
        if (contextId && global.setHexColor) {
            global.setHexColor(contextId, hexColor);
        }
    }
    
    /**
     * Clear the canvas
     */
    clear() {
        const contextId = contexts.get(this._id);
        if (contextId && global.clear) {
            global.clear(contextId);
        }
    }
    
    /**
     * Draw a rectangle
     */
    drawRectangle(x, y, width, height) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawRectangle) {
            global.drawRectangle(contextId, x, y, width, height);
        }
    }
    
    /**
     * Draw a rounded rectangle
     */
    drawRoundedRectangle(x, y, width, height, radius) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawRoundedRectangle) {
            global.drawRoundedRectangle(contextId, x, y, width, height, radius);
        }
    }
    
    /**
     * Draw a circle
     */
    drawCircle(x, y, radius) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawCircle) {
            global.drawCircle(contextId, x, y, radius);
        }
    }
    
    /**
     * Draw an ellipse
     */
    drawEllipse(x, y, rx, ry) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawEllipse) {
            global.drawEllipse(contextId, x, y, rx, ry);
        }
    }
    
    /**
     * Draw a line
     */
    drawLine(x1, y1, x2, y2) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawLine) {
            global.drawLine(contextId, x1, y1, x2, y2);
        }
    }
    
    /**
     * Move to a point
     */
    moveTo(x, y) {
        const contextId = contexts.get(this._id);
        if (contextId && global.moveTo) {
            global.moveTo(contextId, x, y);
        }
    }
    
    /**
     * Draw line to a point
     */
    lineTo(x, y) {
        const contextId = contexts.get(this._id);
        if (contextId && global.lineTo) {
            global.lineTo(contextId, x, y);
        }
    }
    
    /**
     * Draw cubic Bézier curve
     */
    curveTo(cp1x, cp1y, cp2x, cp2y, x, y) {
        const contextId = contexts.get(this._id);
        if (contextId && global.curveTo) {
            global.curveTo(contextId, cp1x, cp1y, cp2x, cp2y, x, y);
        }
    }
    
    /**
     * Close the current path
     */
    closePath() {
        const contextId = contexts.get(this._id);
        if (contextId && global.closePath) {
            global.closePath(contextId);
        }
    }
    
    /**
     * Fill the current path
     */
    fill() {
        const contextId = contexts.get(this._id);
        if (contextId && global.fill) {
            global.fill(contextId);
        }
    }
    
    /**
     * Stroke the current path
     */
    stroke() {
        const contextId = contexts.get(this._id);
        if (contextId && global.stroke) {
            global.stroke(contextId);
        }
    }
    
    /**
     * Set line width
     */
    setLineWidth(width) {
        const contextId = contexts.get(this._id);
        if (contextId && global.setLineWidth) {
            global.setLineWidth(contextId, width);
        }
    }
    
    /**
     * Draw text
     */
    drawString(text, x, y) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawString) {
            global.drawString(contextId, text, x, y);
        }
    }
    
    /**
     * Draw anchored text
     */
    drawStringAnchored(text, x, y, ax, ay) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawStringAnchored) {
            global.drawStringAnchored(contextId, text, x, y, ax, ay);
        }
    }
    
    /**
     * Load font
     */
    loadFont(path, size) {
        const contextId = contexts.get(this._id);
        if (contextId && global.loadFont) {
            global.loadFont(contextId, path, size);
        }
    }
    
    /**
     * Save as PNG
     */
    async savePNG(filePath) {
        const contextId = contexts.get(this._id);
        if (contextId && global.savePNG) {
            return new Promise((resolve, reject) => {
                try {
                    global.savePNG(contextId, filePath);
                    resolve();
                } catch (error) {
                    reject(error);
                }
            });
        }
    }
    
    /**
     * Save as JPEG
     */
    async saveJPEG(filePath, quality = 95) {
        const contextId = contexts.get(this._id);
        if (contextId && global.saveJPEG) {
            return new Promise((resolve, reject) => {
                try {
                    global.saveJPEG(contextId, filePath, quality);
                    resolve();
                } catch (error) {
                    reject(error);
                }
            });
        }
    }
    
    /**
     * Get PNG as Buffer
     */
    async toPNG() {
        const contextId = contexts.get(this._id);
        if (contextId && global.toPNG) {
            return new Promise((resolve, reject) => {
                try {
                    const buffer = global.toPNG(contextId);
                    resolve(Buffer.from(buffer));
                } catch (error) {
                    reject(error);
                }
            });
        }
    }
    
    /**
     * Draw dashed line
     */
    drawDashedLine(x1, y1, x2, y2, pattern) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawDashedLine) {
            global.drawDashedLine(contextId, x1, y1, x2, y2, pattern);
        }
    }
    
    /**
     * Draw text on circle
     */
    drawTextOnCircle(text, x, y, radius) {
        const contextId = contexts.get(this._id);
        if (contextId && global.drawTextOnCircle) {
            global.drawTextOnCircle(contextId, text, x, y, radius);
        }
    }
    
    /**
     * Create linear gradient
     */
    createLinearGradient(x1, y1, x2, y2) {
        if (global.createLinearGradient) {
            const gradientId = global.createLinearGradient(x1, y1, x2, y2);
            return new Gradient(gradientId);
        }
        return null;
    }
    
    /**
     * Set fill style to gradient
     */
    setFillStyle(gradient) {
        const contextId = contexts.get(this._id);
        if (contextId && gradient && global.setFillStyleGradient) {
            global.setFillStyleGradient(contextId, gradient._id);
        }
    }
}

/**
 * Gradient class for advanced fill styles
 */
class Gradient {
    constructor(gradientId) {
        this._id = gradientId;
    }
    
    /**
     * Add color stop to gradient
     */
    addColorStop(position, color) {
        if (global.addColorStop) {
            const [r, g, b, a = 1] = color;
            global.addColorStop(this._id, position, r, g, b, a);
        }
    }
    
    /**
     * Dispose gradient
     */
    dispose() {
        if (global.destroyGradient) {
            global.destroyGradient(this._id);
        }
    }
}

/**
 * Layer Manager for multi-layer compositing
 */
class LayerManager {
    constructor(width, height) {
        ensureInitialized();
        
        this.width = width;
        this.height = height;
        this._layers = new Map();
        
        if (global.createLayerManager) {
            this._id = global.createLayerManager(width, height);
        }
    }
    
    /**
     * Add a new layer
     */
    addLayer(name) {
        if (global.addLayer) {
            const layerContextId = global.addLayer(this._id, name);
            const canvas = new Canvas.__createFromId(layerContextId, this.width, this.height);
            this._layers.set(name, canvas);
            return canvas;
        }
        return null;
    }
    
    /**
     * Set layer opacity
     */
    setLayerOpacity(name, opacity) {
        if (global.setLayerOpacity) {
            global.setLayerOpacity(this._id, name, opacity);
        }
    }
    
    /**
     * Flatten all layers
     */
    flatten() {
        if (global.flattenLayers) {
            const resultContextId = global.flattenLayers(this._id);
            return Canvas.__createFromId(resultContextId, this.width, this.height);
        }
        return null;
    }
    
    /**
     * Dispose layer manager
     */
    dispose() {
        if (global.destroyLayerManager) {
            global.destroyLayerManager(this._id);
        }
        
        // Dispose all layer canvases
        for (const canvas of this._layers.values()) {
            canvas.dispose();
        }
        this._layers.clear();
    }
}

// Helper method for Canvas
Canvas.__createFromId = function(contextId, width, height) {
    const canvas = Object.create(Canvas.prototype);
    canvas.width = width;
    canvas.height = height;
    canvas._id = nextContextId++;
    contexts.set(canvas._id, contextId);
    return canvas;
};

// Utility functions
function hexToRgb(hex) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? [
        parseInt(result[1], 16) / 255,
        parseInt(result[2], 16) / 255,
        parseInt(result[3], 16) / 255
    ] : null;
}

function rgbToHex(r, g, b) {
    return "#" + ((1 << 24) + (Math.round(r * 255) << 16) + (Math.round(g * 255) << 8) + Math.round(b * 255)).toString(16).slice(1);
}

// Export the module
module.exports = {
    init,
    Canvas,
    Gradient,
    LayerManager,
    hexToRgb,
    rgbToHex
};
