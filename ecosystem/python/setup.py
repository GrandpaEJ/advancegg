#!/usr/bin/env python3
"""
AdvanceGG Python Package Setup

High-performance 2D graphics library for Python.
No Go compiler required - includes pre-built native libraries.
"""

import os
import sys
import platform
import subprocess
from pathlib import Path
from setuptools import setup, find_packages, Extension
from setuptools.command.build_ext import build_ext
from setuptools.command.install import install

# Package metadata
PACKAGE_NAME = "advancegg"
VERSION = "1.0.0"
DESCRIPTION = "High-performance 2D graphics library for Python"
LONG_DESCRIPTION = """
# AdvanceGG for Python

Professional 2D graphics library with advanced features:

- **High Performance**: SIMD-optimized operations, memory pooling
- **Rich Features**: Vector graphics, text rendering, image processing
- **Easy to Use**: Pythonic API with comprehensive documentation
- **Cross Platform**: Windows, macOS, Linux support
- **No Dependencies**: Self-contained with pre-built native libraries

## Quick Start

```python
import advancegg

# Create canvas
canvas = advancegg.Canvas(800, 600)

# Set background
canvas.set_rgb(0.1, 0.1, 0.3)
canvas.clear()

# Draw shapes
canvas.set_rgb(1, 0.5, 0)
canvas.draw_circle(400, 300, 100)
canvas.fill()

# Add text
canvas.set_rgb(1, 1, 1)
canvas.draw_string("Hello AdvanceGG!", 300, 350)

# Save image
canvas.save_png("output.png")
```

## Features

- **Drawing**: Shapes, paths, Bézier curves
- **Text**: Unicode, emoji, advanced typography
- **Images**: Loading, processing, filters
- **Colors**: Gradients, color spaces, palettes
- **Layers**: Multi-layer compositing, blend modes
- **Performance**: SIMD, GPU acceleration, profiling

## Documentation

- [API Reference](https://advancegg.dev/docs/python/)
- [Examples](https://github.com/GrandpaEJ/advancegg/tree/main/ecosystem/python/examples)
- [Tutorials](https://advancegg.dev/tutorials/python/)

## Requirements

- Python 3.7+
- No additional dependencies required

## License

MIT License - see LICENSE file for details.
"""

AUTHOR = "AdvanceGG Contributors"
AUTHOR_EMAIL = ""
URL = "https://github.com/GrandpaEJ/advancegg"
LICENSE = "MIT"

CLASSIFIERS = [
    "Development Status :: 5 - Production/Stable",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.7",
    "Programming Language :: Python :: 3.8",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
    "Programming Language :: Python :: 3.12",
    "Topic :: Multimedia :: Graphics",
    "Topic :: Software Development :: Libraries :: Python Modules",
    "Topic :: Scientific/Engineering :: Visualization",
    "Operating System :: OS Independent",
]

KEYWORDS = [
    "graphics", "2d", "drawing", "canvas", "visualization", 
    "charts", "image-processing", "vector-graphics", "typography"
]

# Platform detection
def get_platform_info():
    """Get current platform information for binary selection."""
    system = platform.system().lower()
    machine = platform.machine().lower()
    
    # Normalize architecture names
    arch_map = {
        'x86_64': 'x64',
        'amd64': 'x64',
        'aarch64': 'arm64',
        'arm64': 'arm64',
        'armv7l': 'armv7',
    }
    
    arch = arch_map.get(machine, machine)
    
    return system, arch

# Pre-built binary management
class BinaryManager:
    """Manages pre-built native libraries."""
    
    def __init__(self):
        self.system, self.arch = get_platform_info()
        self.package_dir = Path(__file__).parent
        self.binary_dir = self.package_dir / "advancegg" / "native"
        
    def get_binary_name(self):
        """Get the expected binary name for current platform."""
        if self.system == "windows":
            return f"advancegg-{self.system}-{self.arch}.dll"
        elif self.system == "darwin":
            return f"advancegg-{self.system}-{self.arch}.dylib"
        else:  # Linux and others
            return f"advancegg-{self.system}-{self.arch}.so"
    
    def download_binary(self):
        """Download pre-built binary if not present."""
        binary_name = self.get_binary_name()
        binary_path = self.binary_dir / binary_name
        
        if binary_path.exists():
            print(f"✅ Found native library: {binary_name}")
            return True
        
        # Create directory if it doesn't exist
        self.binary_dir.mkdir(parents=True, exist_ok=True)
        
        # Download URL
        base_url = "https://github.com/GrandpaEJ/advancegg/releases/download/v1.0.0"
        download_url = f"{base_url}/{binary_name}"
        
        print(f"📦 Downloading native library: {binary_name}")
        
        try:
            import urllib.request
            urllib.request.urlretrieve(download_url, binary_path)
            print(f"✅ Downloaded: {binary_name}")
            return True
        except Exception as e:
            print(f"❌ Failed to download {binary_name}: {e}")
            return False
    
    def build_from_source(self):
        """Build native library from source as fallback."""
        print("🔨 Building native library from source...")
        
        try:
            # Check if Go is available
            subprocess.run(["go", "version"], check=True, capture_output=True)
        except (subprocess.CalledProcessError, FileNotFoundError):
            print("❌ Go compiler not found. Please install Go or use pre-built binaries.")
            return False
        
        try:
            # Build the shared library
            build_cmd = [
                "go", "build", 
                "-buildmode=c-shared",
                "-o", str(self.binary_dir / self.get_binary_name()),
                "advancegg_python.go"
            ]
            
            subprocess.run(build_cmd, cwd=self.package_dir, check=True)
            print("✅ Built native library from source")
            return True
        except subprocess.CalledProcessError as e:
            print(f"❌ Failed to build from source: {e}")
            return False

# Custom build command
class CustomBuildExt(build_ext):
    """Custom build command that handles native library."""
    
    def run(self):
        binary_manager = BinaryManager()
        
        # Try to get pre-built binary first
        if not binary_manager.download_binary():
            # Fallback to building from source
            if not binary_manager.build_from_source():
                print("❌ Failed to obtain native library")
                sys.exit(1)
        
        # Continue with normal build
        super().run()

# Custom install command
class CustomInstall(install):
    """Custom install command with post-install verification."""
    
    def run(self):
        super().run()
        
        # Verify installation
        try:
            import advancegg
            print("✅ AdvanceGG installation verified")
        except ImportError as e:
            print(f"❌ Installation verification failed: {e}")

# Package data
def get_package_data():
    """Get package data including native libraries."""
    package_data = {
        'advancegg': [
            'native/*.so',
            'native/*.dll', 
            'native/*.dylib',
            'py.typed',
        ]
    }
    return package_data

# Entry points
def get_entry_points():
    """Get console script entry points."""
    return {
        'console_scripts': [
            'advancegg-info=advancegg.cli:info_command',
            'advancegg-benchmark=advancegg.cli:benchmark_command',
        ]
    }

# Setup configuration
def main():
    setup(
        name=PACKAGE_NAME,
        version=VERSION,
        description=DESCRIPTION,
        long_description=LONG_DESCRIPTION,
        long_description_content_type="text/markdown",
        author=AUTHOR,
        author_email=AUTHOR_EMAIL,
        url=URL,
        license=LICENSE,
        classifiers=CLASSIFIERS,
        keywords=KEYWORDS,
        
        # Package configuration
        packages=find_packages(),
        package_data=get_package_data(),
        include_package_data=True,
        zip_safe=False,
        
        # Requirements
        python_requires=">=3.7",
        install_requires=[
            # No external dependencies - self-contained
        ],
        extras_require={
            'dev': [
                'pytest>=6.0',
                'pytest-cov>=2.0',
                'black>=21.0',
                'flake8>=3.8',
                'mypy>=0.800',
            ],
            'docs': [
                'sphinx>=4.0',
                'sphinx-rtd-theme>=1.0',
                'myst-parser>=0.15',
            ],
            'examples': [
                'numpy>=1.19',
                'matplotlib>=3.3',
                'pillow>=8.0',
            ]
        },
        
        # Entry points
        entry_points=get_entry_points(),
        
        # Custom commands
        cmdclass={
            'build_ext': CustomBuildExt,
            'install': CustomInstall,
        },
        
        # Project URLs
        project_urls={
            'Documentation': 'https://advancegg.dev/docs/python/',
            'Source': 'https://github.com/GrandpaEJ/advancegg',
            'Tracker': 'https://github.com/GrandpaEJ/advancegg/issues',
            'Examples': 'https://github.com/GrandpaEJ/advancegg/tree/main/ecosystem/python/examples',
        },
    )

if __name__ == "__main__":
    main()
