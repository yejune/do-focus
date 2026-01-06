class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do"
  url "https://github.com/yejune/do/archive/refs/tags/v0.1.19.tar.gz"
  sha256 "53c6a08d4ad0697ad7f0090c2eabbca69c831b58ea95401903d74e7fc6ead137"
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
