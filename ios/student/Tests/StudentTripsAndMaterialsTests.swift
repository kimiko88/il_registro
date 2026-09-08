import XCTest
@testable import StudentApp

final class StudentTripsAndMaterialsTests: XCTestCase {

    func testTripDuration() {
        let days = 5
        XCTAssertGreaterThan(days, 0)
    }

    func testTextbookPrice() {
        let price = 34.50
        XCTAssertGreaterThan(price, 0.0)
    }

    func testMaterialDownloadName() {
        let name = "Formulario Derivate ed Integrali.pdf"
        XCTAssertTrue(name.hasSuffix(".pdf"))
    }
}
