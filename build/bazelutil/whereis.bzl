load("@bazel_skylib//lib:shell.bzl", "shell")

def _gen_script_impl(ctx):
    subs = {
        "@@NAME@@": ctx.attr.target.label.name,
        "@@PACKAGE@@": ctx.attr.target.label.package,
        "@@WORKSPACE@@": ctx.attr.target.label.workspace_root,
    }
    out_file = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.expand_template(
        template = ctx.file._template,
        output = out_file,
        substitutions = subs,
    )
    return [
        DefaultInfo(files = depset([out_file])),
    ]

_gen_script = rule(
    implementation = _gen_script_impl,
    attrs = {
        "target": attr.label(
            cfg = "target",
            executable = True,
            mandatory = True,
        ),
        "_template": attr.label(
            default = "@cockroach//build/bazelutil:whereis.sh.in",
            allow_single_file = True,
        ),
    },
)

def whereis_binary(name, target):
    script_name = name + ".sh"
    _gen_script(
        name = script_name,
        target = target,
        testonly = 1,
    )
    native.sh_binary(
        name = name,
        srcs = [script_name],
        data = [
            target,
        ],
        deps = ["@bazel_tools//tools/bash/runfiles"],
        testonly = 1,
    )
