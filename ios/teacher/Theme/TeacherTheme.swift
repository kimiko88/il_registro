import SwiftUI

public struct TeacherTheme {
    public static let primaryGradient = LinearGradient(
        colors: [Color(red: 0.11, green: 0.31, blue: 0.85), Color(red: 0.01, green: 0.52, blue: 0.78)],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )

    public static let primaryColor = Color(red: 0.11, green: 0.31, blue: 0.85)
    public static let cardBackground = Color(UIColor.secondarySystemGroupedBackground)

    // WCAG 2.2 AAA High Contrast
    public static let highContrastPrimary = Color.yellow
    public static let highContrastBackground = Color.black
    public static let highContrastForeground = Color.white
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
    func accessibleTouchTarget() -> some View {
        self.modifier(AccessibleTouchTargetModifier())
    }

    func dyslexiaFriendlyFont() -> some View {
        self.modifier(DyslexiaFriendlyTextModifier(isEnabled: true))
    }
}
