#!/usr/bin/env python3
"""
AdvanceGG Command Line Interface

Provides utility commands for AdvanceGG Python package.
"""

import sys
import time
import platform
import argparse
from pathlib import Path
from typing import Dict, Any

try:
    import advancegg
except ImportError:
    print("❌ AdvanceGG not properly installed")
    sys.exit(1)


def info_command():
    """Display AdvanceGG installation and system information."""
    print("🎨 AdvanceGG Information")
    print("=" * 50)
    
    # Package information
    print(f"📦 Package Version: {advancegg.__version__}")
    print(f"📍 Installation Path: {Path(advancegg.__file__).parent}")
    
    # System information
    print(f"🖥️  Platform: {platform.platform()}")
    print(f"🐍 Python: {platform.python_version()}")
    print(f"🏗️  Architecture: {platform.machine()}")
    
    # Native library information
    try:
        native_info = advancegg.get_native_info()
        print(f"⚡ Native Library: {native_info['version']}")
        print(f"🔧 Build Type: {native_info['build_type']}")
        print(f"📅 Build Date: {native_info['build_date']}")
        
        # Feature support
        features = native_info.get('features', {})
        print("\n🚀 Feature Support:")
        print(f"   SIMD: {'✅' if features.get('simd') else '❌'}")
        print(f"   GPU Acceleration: {'✅' if features.get('gpu') else '❌'}")
        print(f"   Unicode Text: {'✅' if features.get('unicode') else '❌'}")
        print(f"   Image Filters: {'✅' if features.get('filters') else '❌'}")
        
    except Exception as e:
        print(f"⚠️  Native library info unavailable: {e}")
    
    # Performance information
    try:
        perf_info = advancegg.get_performance_info()
        print(f"\n⚡ Performance:")
        print(f"   CPU Cores: {perf_info['cpu_cores']}")
        print(f"   Memory: {perf_info['memory_gb']:.1f} GB")
        if perf_info.get('gpu_name'):
            print(f"   GPU: {perf_info['gpu_name']}")
    except Exception:
        pass


def benchmark_command():
    """Run AdvanceGG performance benchmarks."""
    print("🏃 AdvanceGG Benchmark Suite")
    print("=" * 50)
    
    benchmarks = [
        ("Basic Drawing", benchmark_basic_drawing),
        ("Image Processing", benchmark_image_processing),
        ("Text Rendering", benchmark_text_rendering),
        ("Memory Performance", benchmark_memory),
    ]
    
    results = {}
    
    for name, benchmark_func in benchmarks:
        print(f"\n🔄 Running {name}...")
        try:
            result = benchmark_func()
            results[name] = result
            print(f"✅ {name}: {result['summary']}")
        except Exception as e:
            print(f"❌ {name} failed: {e}")
            results[name] = {"error": str(e)}
    
    # Summary
    print("\n📊 Benchmark Summary")
    print("-" * 30)
    for name, result in results.items():
        if "error" not in result:
            print(f"{name}: {result['summary']}")
        else:
            print(f"{name}: Failed")
    
    # Performance rating
    rating = calculate_performance_rating(results)
    print(f"\n🏆 Overall Performance: {rating}")


def benchmark_basic_drawing() -> Dict[str, Any]:
    """Benchmark basic drawing operations."""
    canvas = advancegg.Canvas(1000, 1000)
    
    # Warm up
    for _ in range(100):
        canvas.draw_circle(500, 500, 50)
        canvas.fill()
    
    # Benchmark
    start_time = time.time()
    operations = 10000
    
    for i in range(operations):
        x = (i % 1000)
        y = (i // 1000) % 1000
        canvas.set_rgb(i / operations, 0.5, 1 - i / operations)
        canvas.draw_circle(x, y, 5)
        canvas.fill()
    
    elapsed = time.time() - start_time
    ops_per_sec = operations / elapsed
    
    return {
        "operations": operations,
        "elapsed": elapsed,
        "ops_per_sec": ops_per_sec,
        "summary": f"{ops_per_sec:.0f} ops/sec"
    }


def benchmark_image_processing() -> Dict[str, Any]:
    """Benchmark image processing operations."""
    # Create test image
    canvas = advancegg.Canvas(512, 512)
    canvas.set_rgb(0.5, 0.5, 0.5)
    canvas.clear()
    
    # Add some content
    for i in range(50):
        canvas.set_rgb(i / 50, 1 - i / 50, 0.5)
        canvas.draw_circle(256 + i * 5, 256 + i * 3, 20)
        canvas.fill()
    
    image = canvas.to_image()
    
    # Benchmark filters
    start_time = time.time()
    operations = 0
    
    # Test various filters
    filters = [
        ("blur", lambda img: advancegg.apply_blur(img, 3.0)),
        ("grayscale", lambda img: advancegg.apply_grayscale(img)),
        ("brightness", lambda img: advancegg.apply_brightness(img, 1.2)),
    ]
    
    for name, filter_func in filters:
        for _ in range(10):
            filtered = filter_func(image)
            operations += 1
    
    elapsed = time.time() - start_time
    ops_per_sec = operations / elapsed
    
    return {
        "operations": operations,
        "elapsed": elapsed,
        "ops_per_sec": ops_per_sec,
        "summary": f"{ops_per_sec:.1f} filters/sec"
    }


def benchmark_text_rendering() -> Dict[str, Any]:
    """Benchmark text rendering operations."""
    canvas = advancegg.Canvas(800, 600)
    
    # Warm up
    for _ in range(50):
        canvas.draw_string("Test", 100, 100)
    
    # Benchmark
    start_time = time.time()
    operations = 1000
    
    texts = [
        "Hello World",
        "AdvanceGG Graphics",
        "Performance Test 🚀",
        "Unicode: 你好世界",
    ]
    
    for i in range(operations):
        text = texts[i % len(texts)]
        x = (i * 10) % 700
        y = (i * 5) % 500 + 50
        canvas.set_rgb(i / operations, 0.5, 1 - i / operations)
        canvas.draw_string(text, x, y)
    
    elapsed = time.time() - start_time
    ops_per_sec = operations / elapsed
    
    return {
        "operations": operations,
        "elapsed": elapsed,
        "ops_per_sec": ops_per_sec,
        "summary": f"{ops_per_sec:.0f} text/sec"
    }


def benchmark_memory() -> Dict[str, Any]:
    """Benchmark memory allocation and cleanup."""
    import gc
    
    # Force garbage collection
    gc.collect()
    
    start_time = time.time()
    canvases = []
    
    # Create many canvases
    for i in range(100):
        canvas = advancegg.Canvas(200, 200)
        canvas.set_rgb(i / 100, 0.5, 1 - i / 100)
        canvas.clear()
        canvases.append(canvas)
    
    # Clean up
    for canvas in canvases:
        canvas.dispose()
    
    gc.collect()
    elapsed = time.time() - start_time
    
    return {
        "canvases": 100,
        "elapsed": elapsed,
        "summary": f"{elapsed:.3f}s for 100 canvases"
    }


def calculate_performance_rating(results: Dict[str, Any]) -> str:
    """Calculate overall performance rating."""
    scores = []
    
    # Basic drawing score
    if "Basic Drawing" in results and "ops_per_sec" in results["Basic Drawing"]:
        drawing_ops = results["Basic Drawing"]["ops_per_sec"]
        if drawing_ops > 50000:
            scores.append(5)
        elif drawing_ops > 30000:
            scores.append(4)
        elif drawing_ops > 15000:
            scores.append(3)
        elif drawing_ops > 5000:
            scores.append(2)
        else:
            scores.append(1)
    
    # Image processing score
    if "Image Processing" in results and "ops_per_sec" in results["Image Processing"]:
        filter_ops = results["Image Processing"]["ops_per_sec"]
        if filter_ops > 100:
            scores.append(5)
        elif filter_ops > 50:
            scores.append(4)
        elif filter_ops > 25:
            scores.append(3)
        elif filter_ops > 10:
            scores.append(2)
        else:
            scores.append(1)
    
    if not scores:
        return "Unknown"
    
    avg_score = sum(scores) / len(scores)
    
    if avg_score >= 4.5:
        return "Excellent ⭐⭐⭐⭐⭐"
    elif avg_score >= 3.5:
        return "Good ⭐⭐⭐⭐"
    elif avg_score >= 2.5:
        return "Average ⭐⭐⭐"
    elif avg_score >= 1.5:
        return "Below Average ⭐⭐"
    else:
        return "Poor ⭐"


def main():
    """Main CLI entry point."""
    parser = argparse.ArgumentParser(
        description="AdvanceGG Python CLI",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # Info command
    info_parser = subparsers.add_parser("info", help="Show installation information")
    
    # Benchmark command
    benchmark_parser = subparsers.add_parser("benchmark", help="Run performance benchmarks")
    
    args = parser.parse_args()
    
    if args.command == "info":
        info_command()
    elif args.command == "benchmark":
        benchmark_command()
    else:
        parser.print_help()


if __name__ == "__main__":
    main()
