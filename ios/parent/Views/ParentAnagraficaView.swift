import SwiftUI

public struct ParentFamilyDataModel: Equatable {
    public var guardianName: String
    public var phone: String
    public var email: String
    public var enrollmentStatus: String

    public init(guardianName: String = "", phone: String = "", email: String = "", enrollmentStatus: String = "") {
        self.guardianName = guardianName
        self.phone = phone
        self.email = email
        self.enrollmentStatus = enrollmentStatus
    }
}

public struct ParentAnagraficaView: View {
    public var familyData: ParentFamilyDataModel

    public init(familyData: ParentFamilyDataModel = ParentFamilyDataModel()) {
        self.familyData = familyData
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if !familyData.guardianName.isEmpty {
                        HStack {
                            Text("Tutore")
                            Spacer()
                            Text(familyData.guardianName)
                                .fontWeight(.bold)
                        }
                    }

                    if !familyData.phone.isEmpty {
                        HStack {
                            Text("Telefono")
                            Spacer()
                            Text(familyData.phone)
                                .fontWeight(.bold)
                                .foregroundColor(.blue)
                        }
                    }

                    if !familyData.email.isEmpty {
                        HStack {
                            Text("Email")
                            Spacer()
                            Text(familyData.email)
                                .foregroundColor(.secondary)
                        }
                    }

                    if !familyData.enrollmentStatus.isEmpty {
                        HStack {
                            Text("Stato Iscrizione")
                            Spacer()
                            Text(familyData.enrollmentStatus)
                                .fontWeight(.bold)
                                .foregroundColor(.green)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}
