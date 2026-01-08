class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.33.tar.gz"
  sha256 "9adce8db2c5e2089f8a8cb26d6024fc04aaa01765a6449f06c669b2d3755603f"
  license "MIT"
  head "https://github.com/yejune/do-focus.git", branch: "main"

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
