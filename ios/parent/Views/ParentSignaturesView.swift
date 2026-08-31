import SwiftUI

struct ParentSignaturesView: View {
    @State private var pdpSigned = false

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Documenti Scolastici da Sottoscrivere")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Piano Didattico Personalizzato (PDP)")
                                .fontWeight(.bold)
                            Spacer()
                            Text(pdpSigned ? "Firmato" : "Richiesto")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(pdpSigned ? Color.green.opacity(0.2) : Color.red.opacity(0.2))
                                .foregroundColor(pdpSigned ? .green : .red)
                                .cornerRadius(4)
                        }
                        Text("Misure compensative e dispensative concordate.")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        if !pdpSigned {
                            Button("Firma Documento") {
                                pdpSigned = true
                            }
                            .buttonStyle(.borderedProminent)
                            .padding(.top, 2)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Firme & Consensi")
        }
    }
}
