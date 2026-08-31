import XCTest
@testable import StudentApp

final class StudentAPIServiceTests: XCTestCase {

    var apiService: HttpStudentAPIService!

    override func setUp() {
        super.setUp()
        apiService = HttpStudentAPIService()
    }

    func testApiServiceInstantiation() {
        XCTAssertNotNil(apiService)
    }
}
