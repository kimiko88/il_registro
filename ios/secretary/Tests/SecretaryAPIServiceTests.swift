import XCTest
@testable import SecretaryApp

final class SecretaryAPIServiceTests: XCTestCase {

    var apiService: HttpSecretaryAPIService!

    override func setUp() {
        super.setUp()
        apiService = HttpSecretaryAPIService()
    }

    func testApiServiceInstantiation() {
        XCTAssertNotNil(apiService)
    }
}
