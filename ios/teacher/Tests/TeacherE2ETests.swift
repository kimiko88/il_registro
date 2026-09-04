import XCTest

final class TeacherE2ETests: XCTestCase {
    var app: XCUIApplication!

    override func setUpWithError() throws {
        try super.setUpWithError()
        continueAfterFailure = false
        guard ProcessInfo.processInfo.environment["XCUI_TARGET_APPLICATION_BUNDLE_PATH"] != nil else {
            throw XCTSkip("Skipping UI test: No target application configured.")
        }
        app = XCUIApplication()
        app.launch()
    }

    override func tearDown() {
        app = nil
        super.tearDown()
    }

    func testTeacherE2EFlow() {
        let emailField = app.textFields["Email Docente"]
        if emailField.exists {
            emailField.tap()
            emailField.typeText("docente@scuola.it")

            let passwordField = app.secureTextFields["Password"]
            passwordField.tap()
            passwordField.typeText("password123")

            let loginButton = app.buttons["ACCEDI AL REGISTRO"]
            loginButton.tap()
        }

        let tabRegistro = app.tabBars.buttons["Registro"]
        XCTAssertTrue(tabRegistro.exists)

        let tabAppello = app.tabBars.buttons["Appello"]
        tabAppello.tap()
        XCTAssertTrue(app.navigationBars["Appello & Presenze"].exists)

        let tabVoti = app.tabBars.buttons["Voti"]
        tabVoti.tap()
        XCTAssertTrue(app.navigationBars["Inserimento Voti"].exists)

        let tabScrutini = app.tabBars.buttons["Scrutini"]
        tabScrutini.tap()
        XCTAssertTrue(app.navigationBars["Gestione Scrutini"].exists)
    }
}
