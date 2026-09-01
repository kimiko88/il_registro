import SwiftUI

public struct SecretaryFascicoloModel: Equatable {
    public let studentNameAndClass: String
    public let fiscalCode: String
    public let sidiCode: String
    public let vaccineStatus: String

    public init(studentNameAndClass: String = "", fiscalCode: String = "", sidiCode: String = "", vaccineStatus: String = "") {
        self.studentNameAndClass = studentNameAndClass
        self.fiscalCode = fiscalCode
        self.sidiCode = sidiCode
        self.vaccineStatus = vaccineStatus
    }
}

public struct SecretaryFascicoloView: View {
    public var fascicolo: SecretaryFascicoloModel

    public init(fascicolo: SecretaryFascicoloModel = SecretaryFascicoloModel()) {
        self.fascicolo = fascicolo
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(fascicolo.studentNameAndClass.isEmpty ? NSLocalizedString("secretary_dashboard_title", comment: "") : fascicolo.studentNameAndClass)) {
                    if !fascicolo.fiscalCode.isEmpty {
                        HStack {
                            Text(NSLocalizedString("profile_title", comment: ""))
                            Spacer()
                            Text(fascicolo.fiscalCode)
                                .fontWeight(.bold)
                        }

                        HStack {
                            Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            Spacer()
                            Text(fascicolo.sidiCode)
                                .fontWeight(.bold)
                                .foregroundColor(.blue)
                        }

                        HStack {
                            Text(NSLocalizedString("classes_title", comment: ""))
                            Spacer()
                            Text(fascicolo.vaccineStatus)
                                .fontWeight(.bold)
                                .foregroundColor(.green)
                        }
                    } else {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
