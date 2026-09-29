{{template "partials/header.tpl" .}}

<div class="search-wrapper">
    <h1>Find Your Next Experience</h1>
    <p>Select a city to discover upcoming Music and Sports events near you.</p>
    
    <div class="search-box-container">
        <input type="text" id="city-search" placeholder="Type a city name (e.g., Toronto)..." autocomplete="off">
        <div id="autocomplete-suggestions" class="suggestions-dropdown"></div>
        <button id="btn-search" disabled>Search</button>
    </div>
</div>

{{template "partials/footer.tpl" .}}
