{{template "partials/header.tpl" .}}

<div class="hero-container">
    <div class="hero-left-content">
        <div class="hero-badge">LESS SCROLLING. MORE GOING.</div>
        <h1>A city of possibilities.<br>Find your next one.</h1>
        <p class="hero-subtitle">Discover music and sports in one place.<br>Choose your city. Find something worth heading out for.</p>
        
        <div class="hero-tags">
            <span>Live music</span>
            <span>Sports & matchdays</span>
            <span>One simple search</span>
        </div>
    </div>
    
    <div class="hero-right-banner">
        <div class="banner-card">
            <div class="banner-top">
                <span class="banner-tag">THE CITY IS CALLING</span>
                <span class="banner-page">01 / 02</span>
            </div>
            <div class="banner-middle">
                <h2>GO<br>OUT.</h2>
            </div>
            <div class="banner-bottom">
                <p>GOOD PLANS.<br>GREAT MEMORIES.</p>
                <span class="banner-arrow">→</span>
            </div>
        </div>
    </div>
</div>

<div class="search-card">
    <h2>Where are we going? <span class="search-hint">Start with a city, then explore what is on.</span></h2>
    
    <form action="/events" method="GET" class="search-form">
        <div class="input-wrapper">
            <span class="search-icon">🔍</span>
            <input type="text" id="city-search" name="city" placeholder="Search a city, e.g. Toronto" autocomplete="off">
            <input type="hidden" id="countryCode" name="countryCode">
            <div id="autocomplete-suggestions" class="suggestions-dropdown"></div>
        </div>
        <button type="submit" id="btn-search" disabled>Explore events &nbsp; →</button>
    </form>
    
    <div class="search-footer">
            <div class="footer-left">
            <span class="search-instructions">Type at least 3 characters and select a suggestion.</span>
            <div id="selection-success-msg" class="success-message">Toronto</div>
        </div>        <span class="search-note">Local sample cities, not Google results</span>
    </div>
</div>

<div class="journey-section">
    <div class="journey-badge">A SMALL PROJECT. THE COMPLETE JOURNEY.</div>
    <h2>From a city to a ticket.</h2>
    <p class="journey-subtitle">Explore the expected experience before building it with Go and Beego.</p>
    
    <div class="steps-grid">
        <div class="step-item">
            <div class="step-num">01</div>
            <h3>Pick a place</h3>
            <p>Find a city with autocomplete.</p>
        </div>
        <div class="step-item">
            <div class="step-num">02</div>
            <h3>Find your event</h3>
            <p>Browse music and sports together.</p>
        </div>
        <div class="step-item">
            <div class="step-num">03</div>
            <h3>View the details</h3>
            <p>Follow the safe demo ticket flow.</p>
        </div>
    </div>
    
    <div class="sample-cities-card">
        <div class="sample-info">
            <strong>Ready-to-test sample cities</strong>
            <p>Toronto and London have sample events. Dhaka demonstrates an empty result.</p>
        </div>
        <div class="sample-buttons">
            <a href="/events?city=Toronto&countryCode=CA" class="btn-sample">Toronto ↗</a>
            <a href="/events?city=London&countryCode=GB" class="btn-sample">London ↗</a>
            <a href="/events?city=Dhaka&countryCode=BD" class="btn-sample">Dhaka · empty ↗</a>
        </div>
    </div>
</div>

{{template "partials/footer.tpl" .}}
