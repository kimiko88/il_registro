import XCTest

final class StudentE2ETests: XCTestCase {
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

    func testStudentE2EFlow() {
        // 1. Enter Email & Password on Login Screen
        let emailField = app.textFields["Email Studente"]
        if emailField.exists {
            emailField.tap()
            emailField.typeText("mario@scuola.it")

            let passwordField = app.secureTextFields["Password"]
            passwordField.tap()
            passwordField.typeText("password123")

            let loginButton = app.buttons["ACCEDI"]
            loginButton.tap()
        }

        // 2. Verify Tab Navigation on Dashboard
        let tabHome = app.tabBars.buttons["Home"]
        XCTAssertTrue(tabHome.exists)

        let tabVoti = app.tabBars.buttons["Voti"]
        tabVoti.tap()
        XCTAssertTrue(app.navigationBars["I Miei Voti"].exists)

        let tabAgenda = app.tabBars.buttons["Agenda"]
        tabAgenda.tap()
        XCTAssertTrue(app.navigationBars["Compiti & Agenda"].exists)

        let tabPresenze = app.tabBars.buttons["Presenze"]
        tabPresenze.tap()
        XCTAssertTrue(app.navigationBars["Presenze & Assenze"].exists)

        let tabPagella = app.tabBars.buttons["Pagella"]
        tabPagella.tap()
        XCTAssertTrue(app.navigationBars["La Mia Pagella"].exists)
    }
}
