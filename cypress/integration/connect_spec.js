describe('Initial Connection', function() {

     Cypress.on("window:before:load", win => {
            cy.spy(win.console, "log", msg => {
                cy.task('log', `console.log --> ${msg}`)
            });
        });

    it('visits the page and takes a screenshot', function() {
        let port = Cypress.env('PORT');
        cy.visit('http://localhost:' + port);
        let filename = 'screenshot';

        cy.screenshot(filename, {
            capture: "viewport"
        })

    });

});