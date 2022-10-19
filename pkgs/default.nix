{
  self,
  buildGoModule,
  lib,
}:
buildGoModule {
  pname = "my-package";
  version = "0.0.0";
  src = self; # + "/src";
  # vendorSha256 should be set to null if dependencies are vendored. If the dependencies aren't
  # vendored, vendorSha256 must be set to a hash of the content of all dependencies. This hash can
  # be found by setting
  # vendorSha256 = lib.fakeSha256;
  # and then running flox build. The build will fail but output the expected sha, which can then be
  # added here.
  vendorSha256 = null;
}
