// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "RegistroScuolaiOS",
    defaultLocalization: "it",
    platforms: [
        .iOS(.v16),
        .macOS(.v13)
    ],
    products: [
        .library(name: "StudentApp", targets: ["StudentApp"]),
        .library(name: "ParentApp", targets: ["ParentApp"]),
        .library(name: "TeacherApp", targets: ["TeacherApp"]),
        .library(name: "SecretaryApp", targets: ["SecretaryApp"])
    ],
    targets: [
        // Student
        .target(
            name: "StudentApp",
            path: "student",
            exclude: ["Tests"],
            resources: [.process("Resources")]
        ),
        .testTarget(
            name: "StudentAppTests",
            dependencies: ["StudentApp"],
            path: "student/Tests"
        ),

        // Parent
        .target(
            name: "ParentApp",
            path: "parent",
            exclude: ["Tests"],
            resources: [.process("Resources")]
        ),
        .testTarget(
            name: "ParentAppTests",
            dependencies: ["ParentApp"],
            path: "parent/Tests"
        ),

        // Teacher
        .target(
            name: "TeacherApp",
            path: "teacher",
            exclude: ["Tests"],
            resources: [.process("Resources")]
        ),
        .testTarget(
            name: "TeacherAppTests",
            dependencies: ["TeacherApp"],
            path: "teacher/Tests"
        ),

        // Secretary
        .target(
            name: "SecretaryApp",
            path: "secretary",
            exclude: ["Tests"],
            resources: [.process("Resources")]
        ),
        .testTarget(
            name: "SecretaryAppTests",
            dependencies: ["SecretaryApp"],
            path: "secretary/Tests"
        )
    ]
)
