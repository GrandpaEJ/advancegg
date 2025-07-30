# AdvanceGG Node.js Integration

Use AdvanceGG from Node.js through WebAssembly (WASM) or native bindings.

## Installation

### Option 1: WebAssembly (Recommended)

```bash
npm install advancegg-wasm
```

### Option 2: Native Bindings (Better Performance)

```bash
# Prerequisites: Go 1.18+, Node.js 16+, node-gyp
npm install advancegg-native
```

## Quick Start

### WebAssembly Version

```javascript
const AdvanceGG = require('advancegg-wasm');

async function main() {
    // Initialize WASM module
    await AdvanceGG.init();
    
    // Create a new canvas
    const canvas = new AdvanceGG.Canvas(800, 600);
    
    // Set background
    canvas.setRGB(0.1, 0.1, 0.3);
    canvas.clear();
    
    // Draw a red circle
    canvas.setRGB(1, 0, 0);
    canvas.drawCircle(400, 300, 100);
    canvas.fill();
    
    // Add text
    canvas.setRGB(1, 1, 1);
    canvas.drawString("Hello from Node.js!", 300, 350);
    
    // Save as PNG
    await canvas.savePNG("output.png");
    
    console.log("Image saved as output.png");
}

main().catch(console.error);
```

### Native Bindings Version

```javascript
const AdvanceGG = require('advancegg-native');

// Create a new canvas
const canvas = new AdvanceGG.Canvas(800, 600);

// Set background
canvas.setRGB(0.1, 0.1, 0.3);
canvas.clear();

// Draw a red circle
canvas.setRGB(1, 0, 0);
canvas.drawCircle(400, 300, 100);
canvas.fill();

// Add text
canvas.setRGB(1, 1, 1);
canvas.drawString("Hello from Node.js!", 300, 350);

// Save as PNG
canvas.savePNG("output.png");
console.log("Image saved as output.png");
```

## Features

### Basic Drawing

```javascript
const canvas = new AdvanceGG.Canvas(800, 600);

// Shapes
canvas.drawRectangle(100, 100, 200, 150);
canvas.drawCircle(400, 300, 80);
canvas.drawEllipse(600, 200, 100, 60);

// Lines and paths
canvas.moveTo(100, 400);
canvas.lineTo(200, 350);
canvas.lineTo(300, 400);
canvas.closePath();
canvas.stroke();
```

### Advanced Features

```javascript
// Gradients
const gradient = canvas.createLinearGradient(0, 0, 200, 0);
gradient.addColorStop(0, [1, 0, 0, 1]); // Red
gradient.addColorStop(1, [0, 0, 1, 1]); // Blue
canvas.setFillStyle(gradient);

// Text effects
canvas.loadFont("arial.ttf", 24);
canvas.drawTextOnCircle("Circular Text", 400, 300, 150);

// Image processing
const image = await canvas.loadImage("input.jpg");
const blurred = canvas.applyBlur(image, 5.0);
canvas.drawImage(blurred, 0, 0);

// Layers
const layerManager = new AdvanceGG.LayerManager(800, 600);
const bgLayer = layerManager.addLayer("background");
const fgLayer = layerManager.addLayer("foreground");

bgLayer.setRGB(0.2, 0.2, 0.4);
bgLayer.clear();

fgLayer.setRGB(1, 1, 0);
fgLayer.drawCircle(400, 300, 100);
fgLayer.fill();

const result = layerManager.flatten();
await result.savePNG("layered.png");
```

### Express.js Integration

```javascript
const express = require('express');
const AdvanceGG = require('advancegg-wasm');

const app = express();

// Initialize AdvanceGG
AdvanceGG.init().then(() => {
    console.log('AdvanceGG initialized');
});

// Generate chart endpoint
app.get('/chart/:type', async (req, res) => {
    try {
        const { type } = req.params;
        const data = req.query.data ? JSON.parse(req.query.data) : [10, 20, 30, 40];
        
        const canvas = await generateChart(type, data);
        const buffer = await canvas.toPNG();
        
        res.set('Content-Type', 'image/png');
        res.send(buffer);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

async function generateChart(type, data) {
    const canvas = new AdvanceGG.Canvas(600, 400);
    
    // Background
    canvas.setRGB(0.95, 0.95, 0.98);
    canvas.clear();
    
    if (type === 'bar') {
        return generateBarChart(canvas, data);
    } else if (type === 'pie') {
        return generatePieChart(canvas, data);
    } else {
        throw new Error('Unsupported chart type');
    }
}

function generateBarChart(canvas, data) {
    const colors = [
        [0.2, 0.4, 0.8], [0.8, 0.2, 0.4], [0.2, 0.8, 0.2], 
        [0.8, 0.5, 0.2], [0.6, 0.2, 0.8]
    ];
    
    const barWidth = 60;
    const spacing = 80;
    const maxValue = Math.max(...data);
    
    data.forEach((value, i) => {
        const x = 100 + i * spacing;
        const height = (value / maxValue) * 250;
        const color = colors[i % colors.length];
        
        // Bar
        canvas.setRGB(...color);
        canvas.drawRoundedRectangle(x, 300 - height, barWidth, height, 8);
        canvas.fill();
        
        // Value label
        canvas.setRGB(0.2, 0.2, 0.2);
        canvas.drawStringAnchored(value.toString(), x + barWidth/2, 280 - height, 0.5, 0.5);
    });
    
    return canvas;
}

app.listen(3000, () => {
    console.log('Server running on http://localhost:3000');
    console.log('Try: http://localhost:3000/chart/bar?data=[25,45,30,60,35]');
});
```

### Real-time Graphics with Socket.io

```javascript
const express = require('express');
const http = require('http');
const socketIo = require('socket.io');
const AdvanceGG = require('advancegg-wasm');

const app = express();
const server = http.createServer(app);
const io = socketIo(server);

// Initialize AdvanceGG
AdvanceGG.init().then(() => {
    console.log('AdvanceGG ready for real-time graphics');
});

io.on('connection', (socket) => {
    console.log('Client connected');
    
    // Handle drawing commands
    socket.on('draw', async (commands) => {
        try {
            const canvas = new AdvanceGG.Canvas(800, 600);
            
            // Execute drawing commands
            for (const cmd of commands) {
                executeDrawCommand(canvas, cmd);
            }
            
            // Convert to base64 and send back
            const buffer = await canvas.toPNG();
            const base64 = buffer.toString('base64');
            
            socket.emit('image', `data:image/png;base64,${base64}`);
        } catch (error) {
            socket.emit('error', error.message);
        }
    });
    
    socket.on('disconnect', () => {
        console.log('Client disconnected');
    });
});

function executeDrawCommand(canvas, cmd) {
    switch (cmd.type) {
        case 'setRGB':
            canvas.setRGB(cmd.r, cmd.g, cmd.b);
            break;
        case 'drawCircle':
            canvas.drawCircle(cmd.x, cmd.y, cmd.radius);
            break;
        case 'fill':
            canvas.fill();
            break;
        case 'drawRectangle':
            canvas.drawRectangle(cmd.x, cmd.y, cmd.width, cmd.height);
            break;
        // Add more commands as needed
    }
}

server.listen(3000, () => {
    console.log('Real-time graphics server running on port 3000');
});
```

## API Reference

### Canvas Class

```javascript
class Canvas {
    constructor(width, height)
    
    // Basic drawing
    drawRectangle(x, y, width, height)
    drawRoundedRectangle(x, y, width, height, radius)
    drawCircle(x, y, radius)
    drawEllipse(x, y, rx, ry)
    drawLine(x1, y1, x2, y2)
    
    // Path operations
    moveTo(x, y)
    lineTo(x, y)
    curveTo(cp1x, cp1y, cp2x, cp2y, x, y)
    closePath()
    
    // Fill and stroke
    fill()
    stroke()
    clear()
    
    // Colors
    setRGB(r, g, b)
    setRGBA(r, g, b, a)
    setHexColor(hexColor)
    
    // Text
    loadFont(path, size)
    drawString(text, x, y)
    drawStringAnchored(text, x, y, ax, ay)
    drawTextOnCircle(text, x, y, radius)
    
    // Images
    async loadImage(path)
    drawImage(image, x, y)
    async savePNG(path)
    async saveJPEG(path, quality = 95)
    async toPNG() // Returns Buffer
    async toJPEG(quality = 95) // Returns Buffer
    
    // Advanced features
    createLinearGradient(x1, y1, x2, y2)
    createRadialGradient(x1, y1, r1, x2, y2, r2)
    setFillStyle(gradient)
    applyBlur(image, radius)
    applyGrayscale(image)
}
```

### LayerManager Class

```javascript
class LayerManager {
    constructor(width, height)
    addLayer(name)
    removeLayer(name)
    setLayerOpacity(name, opacity)
    setLayerBlendMode(name, mode)
    flatten()
}
```

### Gradient Class

```javascript
class Gradient {
    addColorStop(position, color) // color: [r, g, b, a]
}
```

## Performance Optimization

### WebAssembly Tips

```javascript
// Reuse canvas instances
const canvasPool = [];

function getCanvas(width, height) {
    let canvas = canvasPool.find(c => c.width === width && c.height === height);
    if (!canvas) {
        canvas = new AdvanceGG.Canvas(width, height);
        canvasPool.push(canvas);
    }
    return canvas;
}

// Batch operations
function drawMultipleShapes(canvas, shapes) {
    // Group by color to minimize state changes
    const grouped = shapes.reduce((acc, shape) => {
        const key = `${shape.color.r}-${shape.color.g}-${shape.color.b}`;
        if (!acc[key]) acc[key] = [];
        acc[key].push(shape);
        return acc;
    }, {});
    
    Object.entries(grouped).forEach(([colorKey, shapeGroup]) => {
        const [r, g, b] = colorKey.split('-').map(Number);
        canvas.setRGB(r, g, b);
        
        shapeGroup.forEach(shape => {
            canvas.drawCircle(shape.x, shape.y, shape.radius);
        });
        
        canvas.fill();
    });
}
```

### Memory Management

```javascript
// Clean up resources
function cleanup() {
    // Dispose of large canvases when done
    if (largeCanvas) {
        largeCanvas.dispose();
        largeCanvas = null;
    }
    
    // Clear image cache periodically
    if (imageCache.size > 100) {
        imageCache.clear();
    }
}

// Use streams for large images
const fs = require('fs');
const stream = require('stream');

async function processLargeImage(inputPath, outputPath) {
    const canvas = new AdvanceGG.Canvas(4000, 3000);
    
    // Process in chunks to avoid memory issues
    const image = await canvas.loadImage(inputPath);
    const processed = canvas.applyBlur(image, 2.0);
    
    // Stream output to file
    const buffer = await processed.toPNG();
    fs.writeFileSync(outputPath, buffer);
}
```

## Examples

See the `examples/` directory for complete Node.js examples:

- `basic-drawing.js` - Basic shapes and colors
- `express-charts.js` - Web server generating charts
- `real-time-graphics.js` - Socket.io real-time drawing
- `image-processing.js` - Filters and effects
- `data-visualization.js` - Interactive dashboards
- `game-graphics.js` - Sprite generation
- `pdf-generation.js` - PDF reports with graphics

## Troubleshooting

### Common Issues

**Module not found errors**
- Ensure you've installed the correct package (`advancegg-wasm` or `advancegg-native`)
- Check Node.js version compatibility (16+ required)

**WASM initialization fails**
- Make sure to call `await AdvanceGG.init()` before using any functions
- Check browser compatibility for WebAssembly

**Performance issues**
- Use native bindings for CPU-intensive operations
- Implement canvas pooling for frequent operations
- Batch drawing operations when possible

**Memory leaks**
- Dispose of large canvases explicitly
- Clear image caches periodically
- Use streams for large file operations

## Contributing

Contributions to the Node.js integration are welcome! Please see the main AdvanceGG repository for contribution guidelines.

## License

Same as AdvanceGG - MIT License
