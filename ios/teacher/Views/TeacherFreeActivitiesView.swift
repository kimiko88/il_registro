import SwiftUI

struct TeacherFreeActivitiesView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Attività e Ore Funzionali all'Insegnamento")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Ora a Disposizione • Sede Centrale")
                                .fontWeight(.bold)
                            Spacer()
                            Text("1 h")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Mercoledì • 3ª Ora (10:00 - 11:00)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Disponibile per supplenze e sportello")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Ore a Disposizione")
        }
    }
}
