import SwiftUI

struct TeacherGroupsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Laboratori e Classi Articolate")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Laboratorio STEM & Coding")
                                .fontWeight(.bold)
                            Spacer()
                            Text("18 Studenti")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Classi 3ª A / 3ª B • Giovedì 14:30 - 16:30")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Lab Informatica 2 • Registro Dedicato")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Gruppi & Laboratori")
        }
    }
}
