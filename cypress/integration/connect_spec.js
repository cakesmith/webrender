describe('My First Test', function() {

     Cypress.on("window:before:load", win => {
            cy.spy(win.console, "log", msg => {
                cy.task('log', `console.log --> ${msg}`)
            });
        });

    it('visits the page and checks the title', function() {
        let port = Cypress.env('PORT');
        cy.visit('http://localhost:' + port);
    });

});