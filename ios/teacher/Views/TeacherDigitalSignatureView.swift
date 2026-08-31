import SwiftUI

struct TeacherDigitalSignatureView: View {
    @State private var showingOtpAlert = false
    @State private var otpCode = ""
    @State private var isSigned = false

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Firma Elettronica Avanzata (CAD)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text("Verbali di Scrutinio e Registri")
                            .font(.headline)
                        Text("Firma collegiale PAdES con validità legale probatoria.")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }

                Section(header: Text("Documenti da Firmare")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Verbale Scrutinio 3ª A")
                                .fontWeight(.bold)
                            Spacer()
                            Text(isSigned ? "Firmato (PAdES)" : "In Attesa")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(isSigned ? Color.green.opacity(0.2) : Color.orange.opacity(0.2))
                                .foregroundColor(isSigned ? .green : .orange)
                                .cornerRadius(4)
                        }
                        Text("Scrutinio Finale Giugno • 11/12 Docenti Firmati")
                            .font(.subheadline)
                            .foregroundColor(.secondary)

                        if !isSigned {
                            Button(action: { showingOtpAlert = true }) {
                                Label("Firma con Biometria / OTP", systemImage: "signature")
                            }
                            .buttonStyle(.borderedProminent)
                            .padding(.top, 4)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Firma Verbali")
            .sheet(isPresented: $showingOtpAlert) {
                NavigationView {
                    Form {
                        Section(header: Text("Inserimento OTP di Sicurezza")) {
                            TextField("Codice a 6 cifre", text: $otpCode)
                                .keyboardType(.numberPad)
                        }
                    }
                    .navigationTitle("Autorizzazione Firma")
                    .toolbar {
                        ToolbarItem(placement: .cancellationAction) {
                            Button("Annulla") { showingOtpAlert = false }
                        }
                        ToolbarItem(placement: .confirmationAction) {
                            Button("Apponi Firma") {
                                isSigned = true
                                showingOtpAlert = false
                            }
                        }
                    }
                }
            }
        }
    }
}
