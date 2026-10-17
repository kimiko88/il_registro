import SwiftUI

public struct StudentTheme {
    public static let primaryGradient = LinearGradient(
        colors: [Color(red: 0.48, green: 0.22, blue: 0.92), Color(red: 0.31, green: 0.27, blue: 0.90)],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    public static let cardBackground = Color(UIColor.secondarySystemGroupedBackground)
    public static let accentEmerald = Color(red: 0.06, green: 0.72, blue: 0.50)
    public static let accentAmber = Color(red: 0.96, green: 0.62, blue: 0.04)

    // WCAG 2.2 AAA High Contrast
    public static let highContrastPrimary = Color.yellow
    public static let highContrastBackground = Color.black
    public static let highContrastForeground = Color.white
}

public struct GlassCardModifier: ViewModifier {
    public func body(content: Content) -> some View {
        content
            .padding()
            .background(.thinMaterial)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.06), radius: 8, x: 0, y: 4)
    }
}

public struct AccessibleTouchTargetModifier: ViewModifier {
    public func body(content: Content) -> some View {
        content
            .frame(minWidth: 44, minHeight: 44)
    }
}

public struct DyslexiaFriendlyTextModifier: ViewModifier {
    var isEnabled: Bool = true

    public func body(content: Content) -> some View {
        content
            .font(.system(size: 16, design: .rounded))
            .lineSpacing(6)
            .tracking(0.5)
    }
}

public extension View {
    func glassCard() -> some View {
        self.modifier(GlassCardModifier())
    }

    func accessibleTouchTarget() -> some View {
        self.modifier(AccessibleTouchTargetModifier())
    }

    func dyslexiaFriendlyFont() -> some View {
        self.modifier(DyslexiaFriendlyTextModifier(isEnabled: true))
    }
}
