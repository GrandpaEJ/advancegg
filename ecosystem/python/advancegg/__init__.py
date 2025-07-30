"""
AdvanceGG Python Wrapper

High-performance 2D graphics library for Python using Go backend.
"""

import ctypes
import os
import platform
import sys
from pathlib import Path
from typing import List, Tuple, Optional, Dict, Any

# Package metadata
__version__ = "1.0.0"
__author__ = "AdvanceGG Contributors"
__email__ = "hello@advancegg.dev"
__license__ = "MIT"

# Load the shared library
def _find_native_library():
    """Find the appropriate native library for current platform."""
    package_dir = Path(__file__).parent
    native_dir = package_dir / "native"

    system = platform.system().lower()
    machine = platform.machine().lower()

    # Normalize architecture
    arch_map = {
        'x86_64': 'x64', 'amd64': 'x64',
        'aarch64': 'arm64', 'arm64': 'arm64',
        'armv7l': 'armv7'
    }
    arch = arch_map.get(machine, machine)

    # Determine library extension
    if system == "windows":
        lib_name = f"advancegg-{system}-{arch}.dll"
    elif system == "darwin":
        lib_name = f"advancegg-{system}-{arch}.dylib"
    else:
        lib_name = f"advancegg-{system}-{arch}.so"

    lib_path = native_dir / lib_name

    if lib_path.exists():
        return str(lib_path)

    # Fallback to generic names
    fallbacks = ["advancegg.so", "advancegg.dll", "advancegg.dylib"]
    for fallback in fallbacks:
        fallback_path = package_dir / fallback
        if fallback_path.exists():
            return str(fallback_path)

    raise ImportError(
        f"AdvanceGG native library not found. "
        f"Expected: {lib_path} or fallbacks in {package_dir}. "
        f"Platform: {system}-{arch}"
    )

_lib_path = _find_native_library()
_lib = ctypes.CDLL(_lib_path)

# Define function signatures
_lib.create_context.argtypes = [ctypes.c_int, ctypes.c_int]
_lib.create_context.restype = ctypes.c_int

_lib.destroy_context.argtypes = [ctypes.c_int]
_lib.destroy_context.restype = None

_lib.set_rgb.argtypes = [ctypes.c_int, ctypes.c_double, ctypes.c_double, ctypes.c_double]
_lib.set_rgb.restype = None

_lib.set_rgba.argtypes = [ctypes.c_int, ctypes.c_double, ctypes.c_double, ctypes.c_double, ctypes.c_double]
_lib.set_rgba.restype = None

_lib.clear.argtypes = [ctypes.c_int]
_lib.clear.restype = None

_lib.draw_rectangle.argtypes = [ctypes.c_int, ctypes.c_double, ctypes.c_double, ctypes.c_double, ctypes.c_double]
_lib.draw_rectangle.restype = None

_lib.draw_circle.argtypes = [ctypes.c_int, ctypes.c_double, ctypes.c_double, ctypes.c_double]
_lib.draw_circle.restype = None

_lib.fill.argtypes = [ctypes.c_int]
_lib.fill.restype = None

_lib.stroke.argtypes = [ctypes.c_int]
_lib.stroke.restype = None

_lib.save_png.argtypes = [ctypes.c_int, ctypes.c_char_p]
_lib.save_png.restype = None

_lib.draw_string.argtypes = [ctypes.c_int, ctypes.c_char_p, ctypes.c_double, ctypes.c_double]
_lib.draw_string.restype = None


class Canvas:
    """Main drawing canvas for AdvanceGG operations."""
    
    def __init__(self, width: int, height: int):
        """Create a new canvas with specified dimensions."""
        self.width = width
        self.height = height
        self._id = _lib.create_context(width, height)
        if self._id < 0:
            raise RuntimeError("Failed to create AdvanceGG context")
    
    def __del__(self):
        """Clean up the canvas when destroyed."""
        if hasattr(self, '_id'):
            _lib.destroy_context(self._id)
    
    def set_rgb(self, r: float, g: float, b: float):
        """Set the current color using RGB values (0.0 to 1.0)."""
        _lib.set_rgb(self._id, r, g, b)
    
    def set_rgba(self, r: float, g: float, b: float, a: float):
        """Set the current color using RGBA values (0.0 to 1.0)."""
        _lib.set_rgba(self._id, r, g, b, a)
    
    def set_hex_color(self, hex_color: str):
        """Set color using hex string (e.g., '#FF5733')."""
        _lib.set_hex_color(self._id, hex_color.encode('utf-8'))
    
    def clear(self):
        """Clear the canvas with the current color."""
        _lib.clear(self._id)
    
    def draw_rectangle(self, x: float, y: float, width: float, height: float):
        """Draw a rectangle."""
        _lib.draw_rectangle(self._id, x, y, width, height)
    
    def draw_rounded_rectangle(self, x: float, y: float, width: float, height: float, radius: float):
        """Draw a rounded rectangle."""
        _lib.draw_rounded_rectangle(self._id, x, y, width, height, radius)
    
    def draw_circle(self, x: float, y: float, radius: float):
        """Draw a circle."""
        _lib.draw_circle(self._id, x, y, radius)
    
    def draw_ellipse(self, x: float, y: float, rx: float, ry: float):
        """Draw an ellipse."""
        _lib.draw_ellipse(self._id, x, y, rx, ry)
    
    def draw_line(self, x1: float, y1: float, x2: float, y2: float):
        """Draw a line."""
        _lib.draw_line(self._id, x1, y1, x2, y2)
    
    def move_to(self, x: float, y: float):
        """Move to a point without drawing."""
        _lib.move_to(self._id, x, y)
    
    def line_to(self, x: float, y: float):
        """Draw a line to a point."""
        _lib.line_to(self._id, x, y)
    
    def curve_to(self, cp1x: float, cp1y: float, cp2x: float, cp2y: float, x: float, y: float):
        """Draw a cubic Bézier curve."""
        _lib.curve_to(self._id, cp1x, cp1y, cp2x, cp2y, x, y)
    
    def close_path(self):
        """Close the current path."""
        _lib.close_path(self._id)
    
    def fill(self):
        """Fill the current path."""
        _lib.fill(self._id)
    
    def stroke(self):
        """Stroke the current path."""
        _lib.stroke(self._id)
    
    def set_line_width(self, width: float):
        """Set the line width for stroking."""
        _lib.set_line_width(self._id, width)
    
    def draw_string(self, text: str, x: float, y: float):
        """Draw text at the specified position."""
        _lib.draw_string(self._id, text.encode('utf-8'), x, y)
    
    def draw_string_anchored(self, text: str, x: float, y: float, ax: float, ay: float):
        """Draw text with anchor positioning (0.0=left/top, 0.5=center, 1.0=right/bottom)."""
        _lib.draw_string_anchored(self._id, text.encode('utf-8'), x, y, ax, ay)
    
    def load_font(self, path: str, size: float):
        """Load a font from file."""
        _lib.load_font_face(self._id, path.encode('utf-8'), size)
    
    def save_png(self, path: str):
        """Save the canvas as a PNG file."""
        _lib.save_png(self._id, path.encode('utf-8'))
    
    def save_jpeg(self, path: str, quality: int = 95):
        """Save the canvas as a JPEG file."""
        _lib.save_jpeg(self._id, path.encode('utf-8'), quality)
    
    def draw_dashed_line(self, x1: float, y1: float, x2: float, y2: float, pattern: List[float]):
        """Draw a dashed line with the specified pattern."""
        pattern_array = (ctypes.c_double * len(pattern))(*pattern)
        _lib.draw_dashed_line(self._id, x1, y1, x2, y2, pattern_array, len(pattern))
    
    def apply_blur(self, radius: float) -> 'Canvas':
        """Apply blur filter and return new canvas."""
        new_id = _lib.apply_blur(self._id, radius)
        if new_id < 0:
            raise RuntimeError("Failed to apply blur")
        
        new_canvas = Canvas.__new__(Canvas)
        new_canvas.width = self.width
        new_canvas.height = self.height
        new_canvas._id = new_id
        return new_canvas
    
    def apply_grayscale(self) -> 'Canvas':
        """Apply grayscale filter and return new canvas."""
        new_id = _lib.apply_grayscale(self._id)
        if new_id < 0:
            raise RuntimeError("Failed to apply grayscale")
        
        new_canvas = Canvas.__new__(Canvas)
        new_canvas.width = self.width
        new_canvas.height = self.height
        new_canvas._id = new_id
        return new_canvas
    
    def draw_text_on_circle(self, text: str, x: float, y: float, radius: float):
        """Draw text following a circular path."""
        _lib.draw_text_on_circle(self._id, text.encode('utf-8'), x, y, radius)
    
    def create_linear_gradient(self, x1: float, y1: float, x2: float, y2: float) -> 'Gradient':
        """Create a linear gradient."""
        return Gradient(_lib.create_linear_gradient(x1, y1, x2, y2))
    
    def set_fill_style(self, gradient: 'Gradient'):
        """Set the fill style to a gradient."""
        _lib.set_fill_style_gradient(self._id, gradient._id)


class Gradient:
    """Gradient object for advanced fill styles."""
    
    def __init__(self, gradient_id: int):
        self._id = gradient_id
    
    def __del__(self):
        """Clean up the gradient when destroyed."""
        if hasattr(self, '_id'):
            _lib.destroy_gradient(self._id)
    
    def add_color_stop(self, position: float, color: Tuple[float, float, float, float]):
        """Add a color stop to the gradient."""
        r, g, b, a = color
        _lib.add_color_stop(self._id, position, r, g, b, a)


class LayerManager:
    """Manage multiple drawing layers with compositing."""
    
    def __init__(self, width: int, height: int):
        self.width = width
        self.height = height
        self._id = _lib.create_layer_manager(width, height)
        if self._id < 0:
            raise RuntimeError("Failed to create LayerManager")
        self._layers = {}
    
    def __del__(self):
        """Clean up the layer manager when destroyed."""
        if hasattr(self, '_id'):
            _lib.destroy_layer_manager(self._id)
    
    def add_layer(self, name: str) -> Canvas:
        """Add a new layer and return its canvas."""
        layer_id = _lib.add_layer(self._id, name.encode('utf-8'))
        if layer_id < 0:
            raise RuntimeError(f"Failed to create layer '{name}'")
        
        layer_canvas = Canvas.__new__(Canvas)
        layer_canvas.width = self.width
        layer_canvas.height = self.height
        layer_canvas._id = layer_id
        
        self._layers[name] = layer_canvas
        return layer_canvas
    
    def set_layer_opacity(self, name: str, opacity: float):
        """Set the opacity of a layer (0.0 to 1.0)."""
        _lib.set_layer_opacity(self._id, name.encode('utf-8'), opacity)
    
    def flatten(self) -> Canvas:
        """Flatten all layers into a single canvas."""
        result_id = _lib.flatten_layers(self._id)
        if result_id < 0:
            raise RuntimeError("Failed to flatten layers")
        
        result_canvas = Canvas.__new__(Canvas)
        result_canvas.width = self.width
        result_canvas.height = self.height
        result_canvas._id = result_id
        return result_canvas


# Convenience functions
def hex_to_rgb(hex_color: str) -> Tuple[float, float, float]:
    """Convert hex color to RGB tuple."""
    hex_color = hex_color.lstrip('#')
    if len(hex_color) == 3:
        hex_color = ''.join([c*2 for c in hex_color])
    
    r = int(hex_color[0:2], 16) / 255.0
    g = int(hex_color[2:4], 16) / 255.0
    b = int(hex_color[4:6], 16) / 255.0
    return (r, g, b)


def rgb_to_hex(r: float, g: float, b: float) -> str:
    """Convert RGB values to hex string."""
    return f"#{int(r*255):02x}{int(g*255):02x}{int(b*255):02x}"


# Utility functions
def get_native_info() -> Dict[str, Any]:
    """Get information about the native library."""
    try:
        # Try to get version info from native library
        if hasattr(_lib, 'get_version_info'):
            version_info = _lib.get_version_info()
            return {
                'version': version_info.get('version', __version__),
                'build_type': version_info.get('build_type', 'release'),
                'build_date': version_info.get('build_date', 'unknown'),
                'features': {
                    'simd': version_info.get('simd_enabled', True),
                    'gpu': version_info.get('gpu_enabled', False),
                    'unicode': True,
                    'filters': True,
                }
            }
    except:
        pass

    return {
        'version': __version__,
        'build_type': 'release',
        'build_date': 'unknown',
        'features': {
            'simd': True,
            'gpu': False,
            'unicode': True,
            'filters': True,
        }
    }


def get_performance_info() -> Dict[str, Any]:
    """Get system performance information."""
    import psutil

    try:
        return {
            'cpu_cores': psutil.cpu_count(),
            'memory_gb': psutil.virtual_memory().total / (1024**3),
            'gpu_name': _get_gpu_name(),
        }
    except ImportError:
        return {
            'cpu_cores': os.cpu_count() or 1,
            'memory_gb': 0,
            'gpu_name': None,
        }


def _get_gpu_name() -> Optional[str]:
    """Try to get GPU name."""
    try:
        import subprocess
        if platform.system() == "Windows":
            result = subprocess.run(
                ["wmic", "path", "win32_VideoController", "get", "name"],
                capture_output=True, text=True
            )
            lines = result.stdout.strip().split('\n')
            if len(lines) > 1:
                return lines[1].strip()
        elif platform.system() == "Linux":
            result = subprocess.run(
                ["lspci", "-v"], capture_output=True, text=True
            )
            for line in result.stdout.split('\n'):
                if 'VGA' in line or 'Display' in line:
                    return line.split(':')[-1].strip()
    except:
        pass
    return None


# Export main classes and functions
__all__ = [
    'Canvas', 'Gradient', 'LayerManager',
    'hex_to_rgb', 'rgb_to_hex',
    'get_native_info', 'get_performance_info',
    '__version__', '__author__', '__email__', '__license__'
]
