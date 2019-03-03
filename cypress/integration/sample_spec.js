describe('My First Test', function() {
    it('visits the page and checks the title', function() {
        let port = Cypress.env('PORT');
        cy.visit('http://localhost:' + port);
        cy.title().should('equal', 'webrender')
    })
});