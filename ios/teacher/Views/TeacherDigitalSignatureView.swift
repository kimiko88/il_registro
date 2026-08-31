import SwiftUI

public struct DocumentToSignModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let subtitle: String
    public var isSigned: Bool

    public init(id: String = UUID().uuidString, title: String, subtitle: String, isSigned: Bool = false) {
        self.id = id
        self.title = title
        self.subtitle = subtitle
        self.isSigned = isSigned
    }
}

public struct TeacherDigitalSignatureView: View {
    public var documents: [DocumentToSignModel]
    @State private var showingOtpAlert = false
    @State private var otpCode = ""
    @State private var signedDocIds: Set<String> = []

    public init(documents: [DocumentToSignModel] = []) {
        self.documents = documents
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text(NSLocalizedString("sign_hour", comment: ""))
                            .font(.headline)
                        Text(NSLocalizedString("scrutiny_board", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }

                Section(header: Text("Documenti")) {
                    if documents.isEmpty {
                        Text(NSLocalizedString("select_class", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(documents) { doc in
                            let signed = doc.isSigned || signedDocIds.contains(doc.id)
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(doc.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(signed ? "Firmato (PAdES)" : "In Attesa")
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(signed ? Color.green.opacity(0.2) : Color.orange.opacity(0.2))
                                        .foregroundColor(signed ? .green : .orange)
                                        .cornerRadius(4)
                                }
                                Text(doc.subtitle)
                                    .font(.subheadline)
                                    .foregroundColor(.secondary)

                                if !signed {
                                    Button(action: { showingOtpAlert = true }) {
                                        Label(NSLocalizedString("sign_hour", comment: ""), systemImage: "signature")
                                    }
                                    .buttonStyle(.borderedProminent)
                                    .padding(.top, 4)
                                }
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("sign_hour", comment: ""))
            .sheet(isPresented: $showingOtpAlert) {
                NavigationView {
                    Form {
                        Section(header: Text("OTP")) {
                            TextField("OTP", text: $otpCode)
                                .keyboardType(.numberPad)
                        }
                    }
                    .navigationTitle("OTP")
                    .toolbar {
                        ToolbarItem(placement: .cancellationAction) {
                            Button("Annulla") { showingOtpAlert = false }
                        }
                        ToolbarItem(placement: .confirmationAction) {
                            Button("OK") {
                                if let first = documents.first {
                                    signedDocIds.insert(first.id)
                                }
                                showingOtpAlert = false
                            }
                        }
                    }
                }
            }
        }
    }
}
