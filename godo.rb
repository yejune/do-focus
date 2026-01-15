class Godo < Formula
  desc "CLI installer for Do - Claude Code project environment"
  homepage "https://github.com/yejune/do-focus"
  url "https://github.com/yejune/do-focus/archive/refs/tags/v0.1.56.tar.gz"
  sha256 "3352eb7d4a1ddbddd73f13c68e4584b7dce9e065df7b016bc0728eeb6cad325a"
  license "MIT"
  head "https://github.com/yejune/do-focus.git", branch: "main"

  depends_on "go" => :build

  def install
        bin.install "godo"
    bin.install "do-worker"
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
