import XCTest
@testable import TeacherApp

final class TeacherAPIServiceTests: XCTestCase {

    var apiService: HttpTeacherAPIService!

    override func setUp() {
        super.setUp()
        apiService = HttpTeacherAPIService()
    }

    func testApiServiceInstantiation() {
        XCTAssertNotNil(apiService)
    }
}
