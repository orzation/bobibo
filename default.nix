{ lib
, buildGoModule
, fetchFromGitHub
, pkg-config
, nix-update-script
}:
buildGoModule
rec {
  pname = "bobibo";
  version = "1.4.0";

  src = fetchFromGitHub {
    owner = "orzation";
    repo = pname;
    rev = "v${version}";
    hash = "sha256-/esS/CyjHdHMVsIdugRoggecMI//tGuCaayEiNBEocM=";
  };

  vendorHash = "sha256-LzP2pgRheL/NRQmjluqKb8/yxAuFjbmXVU57HdrGSDU=";
  subPackages = [ "cli/" ];

  ldflags = [ "-s" "-w" "-X main.version=${version}" ];

  nativeBuildInputs = [
    pkg-config
  ];

	doCheck = false;

  postInstall = ''
    mv $out/bin/cli $out/bin/bobibo
  '';

  passthru.updateScript = nix-update-script { };

  meta = with lib; {
    description = "A cli-app, convert pictures to ascii arts.";
    homepage = "https://github.com/orzation/bobibo";
    license = licenses.gpl3;
    mainProgram = "bobibo";
    maintainers = with maintainers; [ msqtt ];
  };
}

