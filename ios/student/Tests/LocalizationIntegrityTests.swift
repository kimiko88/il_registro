import XCTest
@testable import StudentApp

final class LocalizationIntegrityTests: XCTestCase {

    func testSupportedLocalesList() {
        let supportedLocales = [
            "it", "en", "de", "fr", "es", "ar", "ro", "ru", "sq", "uk", "zh-Hans"
        ]
        XCTAssertEqual(supportedLocales.count, 11)
        XCTAssertTrue(supportedLocales.contains("it"))
        XCTAssertTrue(supportedLocales.contains("ar"))
        XCTAssertTrue(supportedLocales.contains("zh-Hans"))
    }

    func testDashboardTitleLocalizationResolves() {
        let localizedTitle = studentLocalizedString("dashboard_title")
        XCTAssertNotEqual(localizedTitle, "dashboard_title")
        XCTAssertEqual(localizedTitle, "La Mia Dashboard")
    }

    func testGradesTitleLocalizationResolves() {
        let localizedGrades = studentLocalizedString("grades_title")
        XCTAssertEqual(localizedGrades, "I Miei Voti")
    }

    func testTabBarLocalizationResolvesConciseNames() {
        XCTAssertEqual(studentLocalizedString("tab_home"), "Home")
        XCTAssertEqual(studentLocalizedString("tab_grades"), "Voti")
        XCTAssertEqual(studentLocalizedString("tab_agenda"), "Agenda")
        XCTAssertEqual(studentLocalizedString("tab_attendance"), "Presenze")
        XCTAssertEqual(studentLocalizedString("tab_report_card"), "Pagella")
    }

    func testGradeEvaluationTypeLocalization() {
        let writtenGrade = GradeItemModel(id: "1", subject: "Matematica", grade: 8.0, type: "Written", date: "2026-08-28")
        let oralGrade = GradeItemModel(id: "2", subject: "Storia", grade: 9.0, type: "Oral", date: "2026-08-28")
        let practicalGrade = GradeItemModel(id: "3", subject: "Informatica", grade: 10.0, type: "practical", date: "2026-08-28")

        XCTAssertEqual(writtenGrade.localizedType, "Scritto")
        XCTAssertEqual(oralGrade.localizedType, "Orale")
        XCTAssertEqual(practicalGrade.localizedType, "Pratico")
    }

    func testAttendanceTypeLocalization() {
        let presentRecord = StudentAttendanceRecordModel(id: "1", date: "2026-08-28", type: "present")
        let absentRecord = StudentAttendanceRecordModel(id: "2", date: "2026-08-28", type: "absent")
        let lateRecord = StudentAttendanceRecordModel(id: "3", date: "2026-08-28", type: "late")

        XCTAssertEqual(presentRecord.localizedType, "Presenza")
        XCTAssertEqual(absentRecord.localizedType, "Assenza")
        XCTAssertEqual(lateRecord.localizedType, "Ritardo")
    }
}
