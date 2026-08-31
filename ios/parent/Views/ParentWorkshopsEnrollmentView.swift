import SwiftUI

struct ParentWorkshopsEnrollmentView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Laboratori & Corsi Pomeridiani")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Robotica & Python PNRR")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Iscritto")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Docente: Prof. Neri • Martedì 14:30 - 16:30")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("20/30 Ore Completate • Gratuito PNRR")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Corsi Pomeridiani")
        }
    }
}
