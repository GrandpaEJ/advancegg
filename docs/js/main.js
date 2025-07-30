// AdvanceGG Documentation JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // Initialize components
    initSmoothScrolling();
    initGallery();
    initCodeCopy();
    initNavbarScroll();
    initAnimations();
});

// Smooth scrolling for navigation links
function initSmoothScrolling() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                const offsetTop = target.offsetTop - 80; // Account for fixed navbar
                window.scrollTo({
                    top: offsetTop,
                    behavior: 'smooth'
                });
            }
        });
    });
}

// Initialize gallery
function initGallery() {
    const galleryContainer = document.getElementById('galleryContainer');
    if (!galleryContainer) return;

    // Gallery images data
    const galleryImages = [
        {
            src: 'images/data-visualization-demo.png',
            title: 'Data Visualization',
            description: 'Interactive charts and graphs'
        },
        {
            src: 'images/game-graphics-demo.png',
            title: 'Game Graphics',
            description: 'Sprites and game UI elements'
        },
        {
            src: 'images/creative-text-effects.png',
            title: 'Text Effects',
            description: 'Advanced typography and text-on-path'
        },
        {
            src: 'images/layer-system-demo.png',
            title: 'Layer System',
            description: 'Multi-layer compositing'
        },
        {
            src: 'images/advanced-composite-demo.png',
            title: 'Advanced Compositing',
            description: 'Complex layer blending'
        },
        {
            src: 'images/unicode-shaping-demo.png',
            title: 'Unicode Support',
            description: 'International text rendering'
        },
        {
            src: 'images/emoji-rendering-demo.png',
            title: 'Emoji Rendering',
            description: 'Color emoji support'
        },
        {
            src: 'images/color-space-comparison.png',
            title: 'Color Management',
            description: 'ICC color profiles'
        },
        {
            src: 'images/basic-text-on-path.png',
            title: 'Text on Path',
            description: 'Text following curves'
        },
        {
            src: 'images/advanced-stroke-effects.png',
            title: 'Advanced Strokes',
            description: 'Dashed and gradient strokes'
        },
        {
            src: 'images/filter-showcase.png',
            title: 'Image Filters',
            description: '15+ professional filters'
        },
        {
            src: 'images/performance-demo.png',
            title: 'Performance',
            description: 'Optimized rendering'
        }
    ];

    // Create gallery items
    galleryImages.forEach((image, index) => {
        const galleryItem = createGalleryItem(image, index);
        galleryContainer.appendChild(galleryItem);
    });
}

// Create gallery item
function createGalleryItem(image, index) {
    const col = document.createElement('div');
    col.className = 'col-md-6 col-lg-4 col-xl-3';
    
    col.innerHTML = `
        <div class="gallery-item" data-bs-toggle="modal" data-bs-target="#galleryModal" data-index="${index}">
            <img src="${image.src}" alt="${image.title}" class="img-fluid" loading="lazy" 
                 onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZjhmOWZhIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxNCIgZmlsbD0iIzZjNzU3ZCIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPkltYWdlIG5vdCBmb3VuZDwvdGV4dD48L3N2Zz4='">
            <div class="gallery-overlay">
                <h6>${image.title}</h6>
                <p>${image.description}</p>
            </div>
        </div>
    `;
    
    return col;
}

// Initialize code copy functionality
function initCodeCopy() {
    // Add copy buttons to code blocks
    document.querySelectorAll('pre code').forEach(block => {
        const button = document.createElement('button');
        button.className = 'btn btn-sm btn-outline-light copy-btn';
        button.innerHTML = '<i class="bi bi-clipboard"></i>';
        button.title = 'Copy code';
        
        button.addEventListener('click', () => {
            navigator.clipboard.writeText(block.textContent).then(() => {
                button.innerHTML = '<i class="bi bi-check"></i>';
                button.classList.add('btn-success');
                setTimeout(() => {
                    button.innerHTML = '<i class="bi bi-clipboard"></i>';
                    button.classList.remove('btn-success');
                }, 2000);
            });
        });
        
        const pre = block.parentElement;
        pre.style.position = 'relative';
        pre.appendChild(button);
        
        // Position the button
        button.style.position = 'absolute';
        button.style.top = '10px';
        button.style.right = '10px';
        button.style.zIndex = '10';
    });
}

// Navbar scroll effect
function initNavbarScroll() {
    const navbar = document.querySelector('.navbar');
    let lastScrollTop = 0;
    
    window.addEventListener('scroll', () => {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        
        if (scrollTop > 100) {
            navbar.classList.add('navbar-scrolled');
        } else {
            navbar.classList.remove('navbar-scrolled');
        }
        
        // Hide/show navbar on scroll
        if (scrollTop > lastScrollTop && scrollTop > 200) {
            navbar.style.transform = 'translateY(-100%)';
        } else {
            navbar.style.transform = 'translateY(0)';
        }
        
        lastScrollTop = scrollTop;
    });
}

// Initialize animations
function initAnimations() {
    // Intersection Observer for fade-in animations
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('fade-in-up');
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);
    
    // Observe elements for animation
    document.querySelectorAll('.feature-card, .example-card, .api-section, .step').forEach(el => {
        observer.observe(el);
    });
}

// Gallery modal functionality
function initGalleryModal() {
    // Create modal if it doesn't exist
    if (!document.getElementById('galleryModal')) {
        const modal = document.createElement('div');
        modal.className = 'modal fade';
        modal.id = 'galleryModal';
        modal.innerHTML = `
            <div class="modal-dialog modal-lg modal-dialog-centered">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="galleryModalTitle">Gallery</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body text-center">
                        <img id="galleryModalImage" src="" alt="" class="img-fluid">
                        <p id="galleryModalDescription" class="mt-3"></p>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        <a id="galleryModalLink" href="#" class="btn btn-primary">View Example</a>
                    </div>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    }
    
    // Handle gallery item clicks
    document.addEventListener('click', (e) => {
        const galleryItem = e.target.closest('.gallery-item');
        if (galleryItem) {
            const index = parseInt(galleryItem.dataset.index);
            const image = galleryImages[index];
            
            document.getElementById('galleryModalTitle').textContent = image.title;
            document.getElementById('galleryModalImage').src = image.src;
            document.getElementById('galleryModalImage').alt = image.title;
            document.getElementById('galleryModalDescription').textContent = image.description;
            document.getElementById('galleryModalLink').href = `docs/examples/${image.title.toLowerCase().replace(/\s+/g, '-')}.html`;
        }
    });
}

// Search functionality
function initSearch() {
    const searchInput = document.getElementById('searchInput');
    if (!searchInput) return;
    
    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.toLowerCase();
        const searchResults = document.getElementById('searchResults');
        
        if (query.length < 2) {
            searchResults.innerHTML = '';
            return;
        }
        
        // Simple search implementation
        const results = searchContent(query);
        displaySearchResults(results, searchResults);
    });
}

// Search content
function searchContent(query) {
    const searchableContent = [
        { title: 'Getting Started', url: 'docs/getting-started.html', content: 'installation setup basic usage' },
        { title: 'Drawing API', url: 'docs/api/drawing.html', content: 'shapes circles rectangles lines paths' },
        { title: 'Text Rendering', url: 'docs/api/text.html', content: 'fonts typography unicode emoji' },
        { title: 'Layer System', url: 'docs/api/layers.html', content: 'layers blending compositing opacity' },
        { title: 'Image Processing', url: 'docs/api/images.html', content: 'filters effects blur sharpen' },
        { title: 'Color Management', url: 'docs/api/color.html', content: 'icc profiles color spaces rgb cmyk' }
    ];
    
    return searchableContent.filter(item => 
        item.title.toLowerCase().includes(query) || 
        item.content.toLowerCase().includes(query)
    );
}

// Display search results
function displaySearchResults(results, container) {
    if (results.length === 0) {
        container.innerHTML = '<div class="search-result">No results found</div>';
        return;
    }
    
    container.innerHTML = results.map(result => `
        <div class="search-result">
            <a href="${result.url}" class="search-result-link">
                <h6>${result.title}</h6>
                <p>${result.content}</p>
            </a>
        </div>
    `).join('');
}

// Theme toggle functionality
function initThemeToggle() {
    const themeToggle = document.getElementById('themeToggle');
    if (!themeToggle) return;
    
    const currentTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', currentTheme);
    
    themeToggle.addEventListener('click', () => {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        
        // Update toggle icon
        const icon = themeToggle.querySelector('i');
        icon.className = newTheme === 'dark' ? 'bi bi-sun' : 'bi bi-moon';
    });
}

// Performance monitoring
function initPerformanceMonitoring() {
    // Monitor page load performance
    window.addEventListener('load', () => {
        const perfData = performance.getEntriesByType('navigation')[0];
        console.log('Page load time:', perfData.loadEventEnd - perfData.loadEventStart, 'ms');
    });
    
    // Monitor scroll performance
    let scrollTimeout;
    window.addEventListener('scroll', () => {
        if (scrollTimeout) {
            clearTimeout(scrollTimeout);
        }
        scrollTimeout = setTimeout(() => {
            // Scroll ended
        }, 100);
    });
}

// Error handling
window.addEventListener('error', (e) => {
    console.error('JavaScript error:', e.error);
    // Could send error reports to analytics service
});

// Initialize all functionality when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    initGalleryModal();
    initSearch();
    initThemeToggle();
    initPerformanceMonitoring();
});

// Export functions for use in other scripts
window.AdvanceGGDocs = {
    initSmoothScrolling,
    initGallery,
    initCodeCopy,
    initNavbarScroll,
    initAnimations
};
