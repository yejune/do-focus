class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do"
  url "https://github.com/yejune/do/archive/refs/tags/v0.1.24.tar.gz"
  sha256 "ff962d89d7712fa9f46368edf843cc8abd316c22b0b793b1dcdbd28ed5568f34"
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
