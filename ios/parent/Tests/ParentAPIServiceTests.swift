import XCTest
@testable import ParentApp

final class ParentAPIServiceTests: XCTestCase {

    var apiService: HttpParentAPIService!

    override func setUp() {
        super.setUp()
        apiService = HttpParentAPIService()
    }

    func testApiServiceInstantiation() {
        XCTAssertNotNil(apiService)
    }
}
