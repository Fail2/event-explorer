{{template "partials/header.tpl" .}}

<div class="listing-wrapper">

    <div class="breadcrumb">
        <a href="/">Discover</a> &nbsp;/&nbsp;
        <span class="current-city">{{.City}}</span>
    </div>

    <div class="listing-header-row">
        <div class="header-left-pane">
            <span class="sub-text-green">YOUR CITY. YOUR NEXT PLAN.</span>
            <h1>What is on in {{.City}}.</h1>
            <p class="section-lead-desc">
                Music and sports, loaded together. Find your next reason to go out.
            </p>
        </div>

        <div class="header-right-pane">
            <a href="/" class="btn-change-city">Change city &nbsp;↗</a>
        </div>
    </div>


    <div class="cache-control">
        <form id="cache-form" class="cache-form">
            <input type="text" name="city" value="{{.City}}" class="cache-input" readonly>
            <input type="text" name="countryCode" value="{{.CountryCode}}" class="cache-input" readonly>

            <div class="cache-select-wrapper">
                <select name="category" class="cache-select">
                    <option value="Music" selected>Music</option>
                    <option value="Sports">Sports</option>
                </select>
            </div>

            <button type="submit" class="cache-button cache-button--partial">
                Clear Current City Cache
            </button>

            <a id="clear-whole-cache" class="cache-button cache-button--full">
                Clear Whole Cache
            </a>
        </form>
    </div>


    <!-- Music -->
    <div class="events-category-section">
        <div class="category-header">
            <span class="sub-badge-text">TURN UP THE EVENING</span>

            <div class="title-with-count">
                <h2>Music</h2>
                {{if .MusicEvents}}
                    <span class="count-badge">{{len .MusicEvents}}</span>
                {{end}}
            </div>
        </div>

        {{if .MusicError}}

            <div class="category-error-state">
                <p>⚠️ {{.MusicError}}</p>
            </div>

        {{else if not .MusicEvents}}

            <div class="category-empty-state">
                <p>No upcoming music events found in {{.City}} for the next few weeks.</p>
            </div>

        {{else}}

            <div class="events-cards-grid">
                {{range .MusicEvents}}
                    <a href="/events/{{.ID}}" class="event-display-card">

                        <div class="card-media-wrapper music-bg">
                            {{if .Image}}
                                <img
                                    src="{{.Image}}"
                                    alt="{{.Name}}"
                                    class="event-card-img"
                                >
                            {{else}}
                                <div class="fallback-graphic">
                                    <span class="graphic-waves">|||</span>
                                    <span class="graphic-label">SAMPLE MUSIC EVENT</span>
                                </div>
                            {{end}}
                        </div>

                        <div class="card-body-content">
                            <span class="event-datetime-stamp">{{.Date}}</span>
                            <h3>{{.Name}}</h3>
                            <p class="event-venue-name">{{.Venue}}</p>

                            <div class="card-footer-row">
                                <span class="event-locality-text">{{.City}}</span>
                                <span class="lnk-view-details">
                                    View details &nbsp;↗
                                </span>
                            </div>
                        </div>

                    </a>
                {{end}}
            </div>

        {{end}}
    </div>


    <!-- Sports -->
    <div class="events-category-section">
        <div class="category-header">
            <span class="sub-badge-text">GET INTO THE GAME</span>

            <div class="title-with-count">
                <h2>Sports</h2>
                {{if .SportsEvents}}
                    <span class="count-badge">{{len .SportsEvents}}</span>
                {{end}}
            </div>
        </div>

        {{if .SportsError}}

            <div class="category-error-state">
                <p>⚠️ {{.SportsError}}</p>
            </div>

        {{else if not .SportsEvents}}

            <div class="category-empty-state">
                <p>No live sporting fixtures scheduled in {{.City}} at the moment.</p>
            </div>

        {{else}}

            <div class="events-cards-grid">
                {{range .SportsEvents}}
                    <a href="/events/{{.ID}}" class="event-display-card">

                        <div class="card-media-wrapper sports-bg">
                            {{if .Image}}
                                <img
                                    src="{{.Image}}"
                                    alt="{{.Name}}"
                                    class="event-card-img"
                                >
                            {{else}}
                                <div class="fallback-graphic">
                                    <span class="graphic-court">⊕</span>
                                    <span class="graphic-label">SAMPLE SPORTS EVENT</span>
                                </div>
                            {{end}}
                        </div>

                        <div class="card-body-content">
                            <span class="event-datetime-stamp">{{.Date}}</span>
                            <h3>{{.Name}}</h3>
                            <p class="event-venue-name">{{.Venue}}</p>

                            <div class="card-footer-row">
                                <span class="event-locality-text">{{.City}}</span>
                                <span class="lnk-view-details">
                                    View details &nbsp;↗
                                </span>
                            </div>
                        </div>

                    </a>
                {{end}}
            </div>

        {{end}}
    </div>

</div>

{{template "partials/footer.tpl" .}}

