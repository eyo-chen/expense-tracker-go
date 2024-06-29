const { github, context } = require('@actions/github');

async function run() {
  const token = process.env.GITHUB_TOKEN;
  const octokit = github.getOctokit(token);

  const { data: comments } = await octokit.rest.issues.listComments({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: context.issue.number,
  });

  const botComment = comments.find(comment => comment.user.type === 'Bot' && comment.body.includes('**Total Test Coverage:**'));
  const totalCoverage = process.env.total_coverage;
  const commentBody = `**Total Test Coverage:** ${totalCoverage}%`;

  if (botComment) {
    await octokit.rest.issues.updateComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      comment_id: botComment.id,
      body: commentBody,
    });
  } else {
    await octokit.rest.issues.createComment({
      issue_number: context.issue.number,
      owner: context.repo.owner,
      repo: context.repo.repo,
      body: commentBody,
    });
  }
}

run();
