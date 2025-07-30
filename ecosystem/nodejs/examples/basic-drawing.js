#!/usr/bin/env node

/**
 * Basic Drawing Example for AdvanceGG Node.js
 * 
 * Demonstrates basic shapes, colors, and drawing operations.
 */

const AdvanceGG = require('../index');
const path = require('path');

async function main() {
    try {
        console.log('Initializing AdvanceGG...');
        
        // Initialize the WebAssembly module
        await AdvanceGG.init();
        
        console.log('Creating basic drawing example...');
        
        // Create a new canvas
        const canvas = new AdvanceGG.Canvas(800, 600);
        
        // Set background color (light blue)
        canvas.setRGB(0.9, 0.95, 1.0);
        canvas.clear();
        
        // Draw title
        canvas.setRGB(0.1, 0.1, 0.2);
        canvas.drawStringAnchored("AdvanceGG Node.js - Basic Drawing", 400, 50, 0.5, 0.5);
        
        // Draw rectangles
        canvas.setRGB(0.8, 0.2, 0.2); // Red
        canvas.drawRectangle(50, 100, 150, 100);
        canvas.fill();
        
        canvas.setRGB(0.2, 0.8, 0.2); // Green
        canvas.drawRoundedRectangle(250, 100, 150, 100, 20);
        canvas.fill();
        
        // Draw circles
        canvas.setRGB(0.2, 0.2, 0.8); // Blue
        canvas.drawCircle(125, 300, 60);
        canvas.fill();
        
        canvas.setRGB(0.8, 0.5, 0.2); // Orange
        canvas.drawEllipse(325, 300, 80, 50);
        canvas.fill();
        
        // Draw lines with different styles
        canvas.setRGB(0.4, 0.4, 0.4);
        canvas.setLineWidth(3);
        canvas.drawLine(500, 100, 700, 200);
        canvas.stroke();
        
        // Dashed line
        canvas.setRGB(0.6, 0.2, 0.8);
        canvas.drawDashedLine(500, 250, 700, 350, [10, 5, 2, 5]);
        
        // Draw a path (triangle)
        canvas.setRGB(1.0, 0.8, 0.2); // Yellow
        canvas.moveTo(550, 400);
        canvas.lineTo(650, 400);
        canvas.lineTo(600, 350);
        canvas.closePath();
        canvas.fill();
        
        // Add stroke to triangle
        canvas.setRGB(0.8, 0.6, 0.0);
        canvas.setLineWidth(3);
        canvas.moveTo(550, 400);
        canvas.lineTo(650, 400);
        canvas.lineTo(600, 350);
        canvas.closePath();
        canvas.stroke();
        
        // Draw curved path
        canvas.setRGB(0.8, 0.2, 0.8);
        canvas.setLineWidth(4);
        canvas.moveTo(50, 450);
        canvas.curveTo(150, 400, 250, 500, 350, 450);
        canvas.stroke();
        
        // Add labels
        canvas.setRGB(0.3, 0.3, 0.3);
        canvas.drawStringAnchored("Rectangle", 125, 220, 0.5, 0.5);
        canvas.drawStringAnchored("Rounded Rect", 325, 220, 0.5, 0.5);
        canvas.drawStringAnchored("Circle", 125, 380, 0.5, 0.5);
        canvas.drawStringAnchored("Ellipse", 325, 380, 0.5, 0.5);
        canvas.drawStringAnchored("Line", 600, 130, 0.5, 0.5);
        canvas.drawStringAnchored("Dashed Line", 600, 280, 0.5, 0.5);
        canvas.drawStringAnchored("Triangle", 600, 420, 0.5, 0.5);
        canvas.drawStringAnchored("Bézier Curve", 200, 480, 0.5, 0.5);
        
        // Save the result
        const outputPath = path.join(__dirname, 'nodejs_basic_drawing.png');
        await canvas.savePNG(outputPath);
        
        console.log(`Saved ${outputPath}`);
        
        // Clean up
        canvas.dispose();
        
    } catch (error) {
        console.error('Error:', error.message);
        process.exit(1);
    }
}

// Advanced example with gradients and layers
async function advancedExample() {
    try {
        console.log('Creating advanced example with gradients and layers...');
        
        // Create layer manager
        const layerManager = new AdvanceGG.LayerManager(800, 600);
        
        // Background layer
        const bgLayer = layerManager.addLayer("background");
        const bgGradient = bgLayer.createLinearGradient(0, 0, 0, 600);
        bgGradient.addColorStop(0, [0.95, 0.95, 0.98, 1]);
        bgGradient.addColorStop(1, [0.85, 0.90, 0.95, 1]);
        bgLayer.setFillStyle(bgGradient);
        bgLayer.drawRectangle(0, 0, 800, 600);
        bgLayer.fill();
        
        // Shapes layer
        const shapesLayer = layerManager.addLayer("shapes");
        
        // Draw circles with gradients
        for (let i = 0; i < 5; i++) {
            const x = 150 + i * 120;
            const y = 300;
            const radius = 50;
            
            const gradient = shapesLayer.createLinearGradient(x - radius, y - radius, x + radius, y + radius);
            const hue = i / 5;
            gradient.addColorStop(0, [1, hue, 0.5, 0.8]);
            gradient.addColorStop(1, [0.5, hue, 1, 0.8]);
            
            shapesLayer.setFillStyle(gradient);
            shapesLayer.drawCircle(x, y, radius);
            shapesLayer.fill();
        }
        
        // Text layer
        const textLayer = layerManager.addLayer("text");
        textLayer.setRGB(0.1, 0.1, 0.2);
        textLayer.drawStringAnchored("Advanced AdvanceGG Features", 400, 100, 0.5, 0.5);
        textLayer.drawTextOnCircle("Circular Text Example", 400, 300, 200);
        
        // Set layer opacity
        layerManager.setLayerOpacity("shapes", 0.8);
        
        // Flatten and save
        const result = layerManager.flatten();
        const outputPath = path.join(__dirname, 'nodejs_advanced_example.png');
        await result.savePNG(outputPath);
        
        console.log(`Saved ${outputPath}`);
        
        // Clean up
        layerManager.dispose();
        result.dispose();
        
    } catch (error) {
        console.error('Advanced example error:', error.message);
    }
}

// Chart generation example
async function chartExample() {
    try {
        console.log('Creating chart example...');
        
        const canvas = new AdvanceGG.Canvas(800, 600);
        
        // Background
        canvas.setRGB(0.98, 0.98, 1.0);
        canvas.clear();
        
        // Title
        canvas.setRGB(0.1, 0.1, 0.2);
        canvas.drawStringAnchored("Sales Performance Chart", 400, 50, 0.5, 0.5);
        
        // Sample data
        const data = [
            { label: "Q1", value: 85, color: [0.2, 0.4, 0.8] },
            { label: "Q2", value: 92, color: [0.8, 0.2, 0.4] },
            { label: "Q3", value: 78, color: [0.2, 0.8, 0.2] },
            { label: "Q4", value: 96, color: [0.8, 0.5, 0.2] },
            { label: "Q5", value: 88, color: [0.6, 0.2, 0.8] }
        ];
        
        const barWidth = 80;
        const spacing = 120;
        const maxValue = Math.max(...data.map(d => d.value));
        
        // Draw bars
        data.forEach((item, i) => {
            const x = 150 + i * spacing;
            const height = (item.value / maxValue) * 300;
            const y = 450 - height;
            
            // Bar gradient
            const gradient = canvas.createLinearGradient(x, y, x, y + height);
            gradient.addColorStop(0, [...item.color, 1]);
            gradient.addColorStop(1, [item.color[0] * 0.7, item.color[1] * 0.7, item.color[2] * 0.7, 1]);
            
            canvas.setFillStyle(gradient);
            canvas.drawRoundedRectangle(x, y, barWidth, height, 8);
            canvas.fill();
            
            // Value label
            canvas.setRGB(0.2, 0.2, 0.2);
            canvas.drawStringAnchored(item.value.toString(), x + barWidth/2, y - 20, 0.5, 0.5);
            
            // Category label
            canvas.drawStringAnchored(item.label, x + barWidth/2, 480, 0.5, 0.5);
        });
        
        // Axes
        canvas.setRGB(0.4, 0.4, 0.4);
        canvas.setLineWidth(2);
        canvas.drawLine(100, 450, 700, 450); // X-axis
        canvas.drawLine(100, 450, 100, 150); // Y-axis
        canvas.stroke();
        
        // Save chart
        const outputPath = path.join(__dirname, 'nodejs_chart_example.png');
        await canvas.savePNG(outputPath);
        
        console.log(`Saved ${outputPath}`);
        
        // Clean up
        canvas.dispose();
        
    } catch (error) {
        console.error('Chart example error:', error.message);
    }
}

// Run examples
async function runAllExamples() {
    await main();
    await advancedExample();
    await chartExample();
    
    console.log('All examples completed successfully!');
}

if (require.main === module) {
    runAllExamples().catch(console.error);
}

module.exports = { main, advancedExample, chartExample };
