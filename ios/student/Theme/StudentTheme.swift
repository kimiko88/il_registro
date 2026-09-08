import SwiftUI

struct StudentTheme {
    static let primaryGradient = LinearGradient(
        colors: [Color(red: 0.48, green: 0.22, blue: 0.92), Color(red: 0.31, green: 0.27, blue: 0.90)],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let cardBackground = Color(UIColor.secondarySystemGroupedBackground)
    static let accentEmerald = Color(red: 0.06, green: 0.72, blue: 0.50)
    static let accentAmber = Color(red: 0.96, green: 0.62, blue: 0.04)
}

struct GlassCardModifier: ViewModifier {
    func body(content: Content) -> some View {
        content
            .padding()
            .background(.thinMaterial)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.06), radius: 8, x: 0, y: 4)
    }
}

extension View {
    func glassCard() -> some View {
        self.modifier(GlassCardModifier())
    }
}
