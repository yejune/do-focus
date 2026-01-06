class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do"
  url "https://github.com/yejune/do/archive/refs/tags/v0.1.21.tar.gz"
  sha256 "1cf2e5a1398e3e483f183f27d2ad621340af427c0a2be3a346f50b4257397ccd"
  license "MIT"
  head "https://github.com/yejune/do.git", branch: "main"

  depends_on "go" => :build

  def install
        system "go", "build", "-ldflags", "-X main.version=#{version}", "-o", "godo", "./cmd/godo/"
    bin.install "godo"
  end

  def test
        assert_match "godo version", shell_output("#{bin}/godo version")
  end

  def caveats
    <<~EOS
            godo has been installed!
      
      Usage:
        cd my-project
        godo init      # Install Do
        godo update    # Update Do
    EOS
  end
end
