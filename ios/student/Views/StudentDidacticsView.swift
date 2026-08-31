import SwiftUI

struct StudentDidacticsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Dispense & Materiale Condiviso")) {
                    HStack {
                        VStack(alignment: .leading, spacing: 4) {
                            Text("Formulario Derivate ed Integrali.pdf")
                                .fontWeight(.bold)
                            Text("Matematica • Prof.ssa Bianchi (1.4 MB)")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button(action: {}) {
                            Image(systemName: "arrow.down.circle.fill")
                                .font(.title2)
                                .foregroundColor(.blue)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Materiale Didattico")
        }
    }
}
