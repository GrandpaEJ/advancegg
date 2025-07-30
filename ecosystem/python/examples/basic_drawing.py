#!/usr/bin/env python3
"""
Basic Drawing Example for AdvanceGG Python

Demonstrates basic shapes, colors, and drawing operations.
"""

import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

import advancegg

def main():
    print("Creating basic drawing example...")
    
    # Create a new canvas
    canvas = advancegg.Canvas(800, 600)
    
    # Set background color (light blue)
    canvas.set_rgb(0.9, 0.95, 1.0)
    canvas.clear()
    
    # Draw title
    canvas.set_rgb(0.1, 0.1, 0.2)
    canvas.load_font("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 24)
    canvas.draw_string_anchored("AdvanceGG Python - Basic Drawing", 400, 50, 0.5, 0.5)
    
    # Draw rectangles
    canvas.set_rgb(0.8, 0.2, 0.2)  # Red
    canvas.draw_rectangle(50, 100, 150, 100)
    canvas.fill()
    
    canvas.set_rgb(0.2, 0.8, 0.2)  # Green
    canvas.draw_rounded_rectangle(250, 100, 150, 100, 20)
    canvas.fill()
    
    # Draw circles
    canvas.set_rgb(0.2, 0.2, 0.8)  # Blue
    canvas.draw_circle(125, 300, 60)
    canvas.fill()
    
    canvas.set_rgb(0.8, 0.5, 0.2)  # Orange
    canvas.draw_ellipse(325, 300, 80, 50)
    canvas.fill()
    
    # Draw lines with different styles
    canvas.set_rgb(0.4, 0.4, 0.4)
    canvas.set_line_width(3)
    canvas.draw_line(500, 100, 700, 200)
    canvas.stroke()
    
    # Dashed line
    canvas.set_rgb(0.6, 0.2, 0.8)
    canvas.set_line_width(4)
    canvas.draw_dashed_line(500, 250, 700, 350, [10, 5, 2, 5])
    
    # Draw a path (triangle)
    canvas.set_rgb(1.0, 0.8, 0.2)  # Yellow
    canvas.move_to(550, 400)
    canvas.line_to(650, 400)
    canvas.line_to(600, 350)
    canvas.close_path()
    canvas.fill()
    
    # Add stroke to triangle
    canvas.set_rgb(0.8, 0.6, 0.0)
    canvas.set_line_width(3)
    canvas.move_to(550, 400)
    canvas.line_to(650, 400)
    canvas.line_to(600, 350)
    canvas.close_path()
    canvas.stroke()
    
    # Draw curved path
    canvas.set_rgb(0.8, 0.2, 0.8)
    canvas.set_line_width(4)
    canvas.move_to(50, 450)
    canvas.curve_to(150, 400, 250, 500, 350, 450)
    canvas.stroke()
    
    # Add labels
    canvas.set_rgb(0.3, 0.3, 0.3)
    canvas.load_font("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 14)
    
    canvas.draw_string_anchored("Rectangle", 125, 220, 0.5, 0.5)
    canvas.draw_string_anchored("Rounded Rect", 325, 220, 0.5, 0.5)
    canvas.draw_string_anchored("Circle", 125, 380, 0.5, 0.5)
    canvas.draw_string_anchored("Ellipse", 325, 380, 0.5, 0.5)
    canvas.draw_string_anchored("Line", 600, 130, 0.5, 0.5)
    canvas.draw_string_anchored("Dashed Line", 600, 280, 0.5, 0.5)
    canvas.draw_string_anchored("Triangle", 600, 420, 0.5, 0.5)
    canvas.draw_string_anchored("Bézier Curve", 200, 480, 0.5, 0.5)
    
    # Save the result
    canvas.save_png("python_basic_drawing.png")
    print("Saved python_basic_drawing.png")

if __name__ == "__main__":
    main()
