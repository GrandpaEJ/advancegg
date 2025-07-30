/**
 * AdvanceGG TypeScript Definitions
 * 
 * High-performance 2D graphics library for Node.js
 */

declare module 'advancegg' {
    /**
     * Initialize the AdvanceGG library
     */
    export function init(): Promise<void>;

    /**
     * Main drawing canvas
     */
    export class Canvas {
        readonly width: number;
        readonly height: number;

        constructor(width: number, height: number);

        // Color operations
        setRGB(r: number, g: number, b: number): void;
        setRGBA(r: number, g: number, b: number, a: number): void;
        setHexColor(hexColor: string): void;

        // Basic drawing
        clear(): void;
        drawRectangle(x: number, y: number, width: number, height: number): void;
        drawRoundedRectangle(x: number, y: number, width: number, height: number, radius: number): void;
        drawCircle(x: number, y: number, radius: number): void;
        drawEllipse(x: number, y: number, rx: number, ry: number): void;
        drawLine(x1: number, y1: number, x2: number, y2: number): void;

        // Path operations
        moveTo(x: number, y: number): void;
        lineTo(x: number, y: number): void;
        curveTo(cp1x: number, cp1y: number, cp2x: number, cp2y: number, x: number, y: number): void;
        closePath(): void;

        // Fill and stroke
        fill(): void;
        stroke(): void;
        setLineWidth(width: number): void;

        // Text rendering
        drawString(text: string, x: number, y: number): void;
        drawStringAnchored(text: string, x: number, y: number, ax: number, ay: number): void;
        loadFont(path: string, size: number): void;
        drawTextOnCircle(text: string, x: number, y: number, radius: number): void;

        // Advanced drawing
        drawDashedLine(x1: number, y1: number, x2: number, y2: number, pattern: number[]): void;
        createLinearGradient(x1: number, y1: number, x2: number, y2: number): Gradient;
        setFillStyle(gradient: Gradient): void;

        // Image operations
        loadImage(path: string): Promise<Image>;
        drawImage(image: Image, x: number, y: number): void;

        // File operations
        savePNG(path: string): Promise<void>;
        saveJPEG(path: string, quality?: number): Promise<void>;
        toPNG(): Promise<Buffer>;
        toJPEG(quality?: number): Promise<Buffer>;

        // Resource management
        dispose(): void;
    }

    /**
     * Gradient for advanced fill styles
     */
    export class Gradient {
        addColorStop(position: number, color: [number, number, number, number]): void;
        dispose(): void;
    }

    /**
     * Image object
     */
    export class Image {
        readonly width: number;
        readonly height: number;
        dispose(): void;
    }

    /**
     * Layer manager for multi-layer compositing
     */
    export class LayerManager {
        readonly width: number;
        readonly height: number;

        constructor(width: number, height: number);

        addLayer(name: string): Canvas;
        removeLayer(name: string): void;
        setLayerOpacity(name: string, opacity: number): void;
        setLayerVisible(name: string, visible: boolean): void;
        setLayerBlendMode(name: string, mode: BlendMode): void;
        flatten(): Canvas;
        dispose(): void;
    }

    /**
     * Blend modes for layer compositing
     */
    export enum BlendMode {
        Normal = 'normal',
        Multiply = 'multiply',
        Screen = 'screen',
        Overlay = 'overlay',
        SoftLight = 'soft-light',
        HardLight = 'hard-light',
        ColorDodge = 'color-dodge',
        ColorBurn = 'color-burn',
        Darken = 'darken',
        Lighten = 'lighten',
        Difference = 'difference',
        Exclusion = 'exclusion'
    }

    /**
     * Utility functions
     */
    export function hexToRgb(hex: string): [number, number, number] | null;
    export function rgbToHex(r: number, g: number, b: number): string;

    /**
     * System information
     */
    export interface SystemInfo {
        platform: string;
        arch: string;
        nodeVersion: string;
        packageVersion: string;
        nativeLibrary: {
            path: string;
            version: string;
            features: {
                simd: boolean;
                gpu: boolean;
                unicode: boolean;
                filters: boolean;
            };
        };
    }

    export function getSystemInfo(): SystemInfo;

    /**
     * Performance monitoring
     */
    export interface PerformanceMetrics {
        fps: number;
        memoryUsage: number;
        cpuUsage: number;
        gpuUsage?: number;
    }

    export function getPerformanceMetrics(): PerformanceMetrics;

    /**
     * Image processing filters
     */
    export function applyBlur(image: Image, radius: number): Image;
    export function applyGrayscale(image: Image): Image;
    export function applyBrightness(image: Image, factor: number): Image;
    export function applyContrast(image: Image, factor: number): Image;
    export function applySaturation(image: Image, factor: number): Image;

    /**
     * Batch processing
     */
    export function processImagesParallel<T>(
        images: Image[],
        processor: (image: Image) => T
    ): Promise<T[]>;

    /**
     * Configuration
     */
    export interface Config {
        simdEnabled: boolean;
        memoryPoolingEnabled: boolean;
        cachingEnabled: boolean;
        workerThreads: number;
        gpuAcceleration: boolean;
    }

    export function setConfig(config: Partial<Config>): void;
    export function getConfig(): Config;

    /**
     * Error types
     */
    export class AdvanceGGError extends Error {
        readonly code: string;
        constructor(message: string, code: string);
    }

    export class InitializationError extends AdvanceGGError {
        constructor(message: string);
    }

    export class NativeLibraryError extends AdvanceGGError {
        constructor(message: string);
    }

    export class RenderingError extends AdvanceGGError {
        constructor(message: string);
    }

    /**
     * Events
     */
    export interface EventEmitter {
        on(event: string, listener: (...args: any[]) => void): this;
        off(event: string, listener: (...args: any[]) => void): this;
        emit(event: string, ...args: any[]): boolean;
    }

    export const events: EventEmitter;

    /**
     * Version information
     */
    export const version: string;
    export const buildDate: string;
    export const gitCommit: string;

    /**
     * Feature detection
     */
    export function isFeatureSupported(feature: string): boolean;
    export function getSupportedFeatures(): string[];

    /**
     * Debugging and profiling
     */
    export interface ProfileResult {
        name: string;
        duration: number;
        percentage: number;
        children?: ProfileResult[];
    }

    export class Profiler {
        start(): void;
        stop(): void;
        beginSection(name: string): void;
        endSection(name: string): void;
        getResults(): ProfileResult[];
        reset(): void;
    }

    export function createProfiler(): Profiler;

    /**
     * Memory management
     */
    export interface MemoryStats {
        totalAllocated: number;
        currentUsage: number;
        peakUsage: number;
        poolStats: {
            [poolName: string]: {
                size: number;
                used: number;
                available: number;
            };
        };
    }

    export function getMemoryStats(): MemoryStats;
    export function forceGarbageCollection(): void;
    export function enableMemoryPooling(enabled: boolean): void;

    /**
     * Logging
     */
    export enum LogLevel {
        Error = 0,
        Warn = 1,
        Info = 2,
        Debug = 3,
        Trace = 4
    }

    export function setLogLevel(level: LogLevel): void;
    export function log(level: LogLevel, message: string, ...args: any[]): void;
}
