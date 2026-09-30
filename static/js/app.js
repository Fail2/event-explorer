document.addEventListener('DOMContentLoaded', function () {
    const cityInput = document.getElementById('city-search');
    const suggestionsDropdown = document.getElementById('autocomplete-suggestions');
    const searchButton = document.getElementById('btn-search');
    const countryCodeInput = document.getElementById('countryCode');

    let sessionToken = generateUUID();

    function generateUUID() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    cityInput.addEventListener('input', function () {
        const query = cityInput.value.trim();

        if (query.length < 3) {
            suggestionsDropdown.innerHTML = '';
            suggestionsDropdown.style.display = 'none';
            searchButton.disabled = true;
            return;
        }
        console.log(query)

        fetch(`/api/locations/autocomplete?input=${encodeURIComponent(query)}&sessionToken=${sessionToken}`)
            .then(response => response.json())
            .then(data => {
                suggestionsDropdown.innerHTML = '';
                console.log(data)

                if (data.suggestions && data.suggestions.length > 0) {
                    data.suggestions.forEach(item => {
                        const div = document.createElement('div');
                        div.className = 'suggestion-item';
                        div.textContent = item.text;
                        div.dataset.placeId = item.placeId;

                        div.addEventListener('click', function () {
                            cityInput.value = item.text;
                            suggestionsDropdown.innerHTML = '';
                            suggestionsDropdown.style.display = 'none';

                            fetch(`/api/locations/${item.placeId}?sessionToken=${sessionToken}`)
                                .then(res => res.json())
                                .then(locationData => {
                                    if (locationData.city && locationData.countryCode) {
                                        cityInput.value = locationData.city;
                                        countryCodeInput.value = locationData.countryCode;
                                        searchButton.disabled = false;
                                        sessionToken = generateUUID();
                                    }
                                    console.log(locationData)
                                    const successMsg = document.querySelector('.success-message');
                                    if (successMsg) {
                                        successMsg.textContent = `Selected ${locationData.city}. You are ready to explore.`;
                                        successMsg.style.display = "block";
                                    }
                                });
                        });

                        suggestionsDropdown.appendChild(div);
                    });
                    suggestionsDropdown.style.display = 'block';
                } else {
                    suggestionsDropdown.style.display = 'none';
                }
            })
            .catch(() => {
                suggestionsDropdown.style.display = 'none';
            });
    });

    document.addEventListener('click', function (e) {
        if (e.target !== cityInput && e.target !== suggestionsDropdown) {
            suggestionsDropdown.style.display = 'none';
        }
    });
});
