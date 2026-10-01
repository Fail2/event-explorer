document.addEventListener('DOMContentLoaded', function () {
    const cacheForm = document.getElementById('cache-form');
    const clearWholeCache = document.getElementById('clear-whole-cache');

    if (cacheForm) {
        cacheForm.addEventListener('submit', async function (event) {
            event.preventDefault();

            const formData = new FormData(cacheForm);

            try {
                const response = await fetch('/api/cache/clear', {
                    method: 'POST',
                    body: formData,
                    headers: {
                        'Accept': 'application/json'
                    }
                });

                const data = await response.json();

                console.log(data);
            } catch (error) {
                console.error('Cache clear failed:', error);
            }
        });
    }

    if (clearWholeCache) {
        clearWholeCache.addEventListener('click', async function (event) {
            event.preventDefault();

            try {
                const response = await fetch('/api/cache/clear', {
                    method: 'POST',
                    headers: {
                        'Accept': 'application/json'
                    }
                });

                const data = await response.json();

                console.log(data);
            } catch (error) {
                console.error('Cache clear failed:', error);
            }
        });
    }
});