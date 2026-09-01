import XCTest

final class ParentE2ETests: XCTestCase {
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

    func testParentE2EFlow() {
        let emailField = app.textFields["Email Genitore"]
        if emailField.exists {
            emailField.tap()
            emailField.typeText("genitore@famiglia.it")

            let passwordField = app.secureTextFields["Password"]
            passwordField.tap()
            passwordField.typeText("password123")

            let loginButton = app.buttons["ACCEDI AL PORTALE"]
            loginButton.tap()
        }

        let tabFigli = app.tabBars.buttons["Figli"]
        XCTAssertTrue(tabFigli.exists)

        let tabVoti = app.tabBars.buttons["Voti"]
        tabVoti.tap()
        XCTAssertTrue(app.navigationBars["Voti Figlio"].exists)

        let tabGiustifiche = app.tabBars.buttons["Giustifiche"]
        tabGiustifiche.tap()
        XCTAssertTrue(app.navigationBars["Giustifiche"].exists)

        let tabColloqui = app.tabBars.buttons["Colloqui"]
        tabColloqui.tap()
        XCTAssertTrue(app.navigationBars["Ricevimento Famiglie"].exists)
    }
}
