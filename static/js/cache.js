document.addEventListener('DOMContentLoaded', function () {
    const cacheForm = document.getElementById('cache-form');
    const clearWholeCache = document.getElementById('clear-whole-cache');

    async function clearCache(body = null) {
        const response = await fetch('/api/cache/clear', {
            method: 'POST',
            body: body
        });

        const text = await response.text();

        console.log('CACHE RESPONSE:', response.status, text);

        if (!response.ok) {
            throw new Error(`Cache clear failed: ${response.status}`);
        }

        window.location.reload();
    }

    if (cacheForm) {
        cacheForm.addEventListener('submit', function (event) {
            event.preventDefault();

            clearCache(new FormData(cacheForm))
                .catch(error => console.error('CACHE ERROR:', error));
        });
    }

    if (clearWholeCache) {
        clearWholeCache.addEventListener('click', function (event) {
            event.preventDefault();

            clearCache()
                .catch(error => console.error('CACHE ERROR:', error));
        });
    }
});