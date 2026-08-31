import SwiftUI

struct ParentTransportView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Linea Scuolabus Assegnata")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Linea 3: Nord - Centro")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Regolare")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Andata: Via Roma, 45 (07:35)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Ritorno: Via Roma, 45 (13:40)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Autista: Sig. Franco • Bus FX920KL")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Scuolabus")
        }
    }
}
