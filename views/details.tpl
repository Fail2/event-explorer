{{template "partials/header.tpl" .}}

<div class="details-wrapper">
    <div class="breadcrumb">
        <a href="/">Discover</a> &nbsp;/&nbsp; 
        <a href="/events?city={{.Event.City}}&countryCode={{.Event.Country}}">{{.Event.City}}</a> &nbsp;/&nbsp; 
        <span class="current-event">Event details</span>
    </div>

    <a href="/events?city={{.Event.City}}&countryCode={{.Event.Country}}" class="btn-back-navigation-top">
        ← Back to events
    </a>

    <div class="details-container-grid">
        <div class="details-left-pane">
            <div class="details-image-hero">
                {{if .Event.Image}}
                    <img src="{{.Event.Image}}" alt="{{.Event.Name}}" class="hero-display-img">
                {{else}}
                    <div class="details-fallback-graphic">
                        <span class="graphic-waves">|||</span>
                        <span class="graphic-label-badge">FICTIONAL SAMPLE EVENT</span>
                    </div>
                {{end}}
            </div>

            <div class="details-info-section-left">
                <span class="details-category-sub-tag">{{.Event.Category}}</span>
                <h1 class="details-main-title-text">{{.Event.Name}}</h1>
                
                <div class="details-about-block">
                    <h2>About this event</h2>
                    <p class="desc-p1">
                        {{if .Event.Description}}
                            {{.Event.Description}}
                        {{else}}
                            This is a fictional event created for the Event Explorer internship preview. It is not a real event and there are no tickets for sale.
                        {{end}}
                    </p>
                    <p class="desc-p2">
                        Use this page to explore the expected details layout and ticket journey. The real assignment renders this information with Beego templates and retrieves it from Ticketmaster.
                    </p>
                </div>
            </div>
        </div>

        <div class="details-right-pane">
            <div class="sticky-sidebar-card">
                <span class="sidebar-badge-accent">MAKE A PLAN</span>
                <h2 class="sidebar-event-title">The details</h2>
                
                <div class="sidebar-meta-item-block">
                    <span class="meta-label">WHEN</span>
                    <strong>{{.Event.Date}}</strong>
                    <p>{{.Event.Time}} ({{.Event.TimeZone}})</p>
                </div>

                <div class="sidebar-meta-item-block">
                    <span class="meta-label">WHERE</span>
                    <strong>{{.Event.Venue}}</strong>
                    <p>{{.Event.City}}, {{.Event.Country}}</p>
                </div>

                <div class="sidebar-meta-item-block">
                    <span class="meta-label">CATEGORY</span>
                    <p class="category-plain-text">{{.Event.Category}}</p>
                </div>

                <a href="{{.Event.TicketURL}}" target="_blank" rel="noopener noreferrer" class="btn-checkout-tickets">
                    <span>View tickets</span>
                    <span class="btn-arrow-icon">↗</span>
                </a>
                
                <div class="sidebar-footer-note">
                    Opens a safe, local ticket-flow simulation. No booking, payment, or real provider visit takes place.
                </div>
                
                <div class="sidebar-bottom-sub-text">
                    SAMPLE DATA / LOCAL DEMO
                </div>
            </div>
        </div>
    </div>
</div>

{{template "partials/footer.tpl" .}}
