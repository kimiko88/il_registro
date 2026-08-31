import XCTest

final class SecretaryE2ETests: XCTestCase {
    var app: XCUIApplication!

    override func setUp() {
        super.setUp()
        continueAfterFailure = false
        app = XCUIApplication()
        app.launch()
    }

    override func tearDown() {
        app = nil
        super.tearDown()
    }

    func testSecretaryE2EFlow() {
        let emailField = app.textFields["Email Segreteria"]
        if emailField.exists {
            emailField.tap()
            emailField.typeText("segreteria@scuola.it")

            let passwordField = app.secureTextFields["Password"]
            passwordField.tap()
            passwordField.typeText("password123")

            let loginButton = app.buttons["ACCEDI AL PANNELLO"]
            loginButton.tap()
        }

        let tabPannello = app.tabBars.buttons["Pannello"]
        XCTAssertTrue(tabPannello.exists)

        let tabUtenti = app.tabBars.buttons["Utenti"]
        tabUtenti.tap()
        XCTAssertTrue(app.navigationBars["Anagrafica Utenti"].exists)

        let tabScrutini = app.tabBars.buttons["Scrutini"]
        tabScrutini.tap()
        XCTAssertTrue(app.navigationBars["Supervisione Scrutini"].exists)

        let tabCertificati = app.tabBars.buttons["Certificati"]
        tabCertificati.tap()
        XCTAssertTrue(app.navigationBars["Emissione Certificati"].exists)

        let tabAudit = app.tabBars.buttons["Audit"]
        tabAudit.tap()
        XCTAssertTrue(app.navigationBars["Registro Audit"].exists)
    }
}
